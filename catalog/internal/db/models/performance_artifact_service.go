package models

import (
	"cmp"
	"fmt"
	"math"
	"slices"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/scopes"
)

type PerformanceArtifactParams struct {
	ModelID               int32
	TargetRPS             int32
	Recommendations       bool
	FilterQuery           string
	PageSize              int32
	OrderBy               string
	SortOrder             string
	NextPageToken         *string
	RPSProperty           string // configurable "requests_per_second"
	LatencyProperty       string // configurable "ttft_p90"
	HardwareCountProperty string // configurable "hardware_count"
	HardwareTypeProperty  string // configurable "hardware_type"
}

type PerformanceArtifactService struct {
	artifactRepo CatalogArtifactRepository
}

func NewPerformanceArtifactService(repo CatalogArtifactRepository) *PerformanceArtifactService {
	return &PerformanceArtifactService{
		artifactRepo: repo,
	}
}

func (s *PerformanceArtifactService) GetArtifacts(params PerformanceArtifactParams) (*models.ListWrapper[CatalogMetricsArtifact], error) {
	// Build filter query to include only performance-metrics
	filterQuery := s.buildPerformanceFilterQuery(params.FilterQuery)

	// Configure repository options
	listOptions := CatalogArtifactListOptions{
		ParentResourceID:    &params.ModelID,
		ArtifactTypesFilter: []string{"metrics-artifact"},
		Pagination: models.Pagination{
			FilterQuery:   &filterQuery,
			PageSize:      &params.PageSize,
			OrderBy:       &params.OrderBy,
			SortOrder:     &params.SortOrder,
			NextPageToken: params.NextPageToken,
		},
	}

	// Recommendations need to be based on the full list. Pagination is handled below.
	if params.Recommendations {
		listOptions.PageSize = nil
		listOptions.NextPageToken = nil
	}

	// Get artifacts from repository
	dbResult, err := s.artifactRepo.List(listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}

	artifacts := make([]CatalogMetricsArtifact, len(dbResult.Items))
	for i := range dbResult.Items {
		artifacts[i] = dbResult.Items[i].CatalogMetricsArtifact
	}

	// Validate custom properties exist
	if err := s.validateCustomProperties(artifacts, params.RPSProperty, params.LatencyProperty, params.HardwareCountProperty, params.HardwareTypeProperty); err != nil {
		return nil, fmt.Errorf("invalid custom properties: %w", err)
	}

	// Apply performance-specific processing
	artifacts = s.processArtifacts(artifacts, params)

	list := &models.ListWrapper[CatalogMetricsArtifact]{
		Items:         artifacts,
		PageSize:      dbResult.PageSize,
		Size:          dbResult.Size,
		NextPageToken: dbResult.NextPageToken,
	}

	if params.Recommendations {
		s.paginate(list, params.PageSize, params.NextPageToken)
	}

	return list, nil
}

// buildPerformanceFilterQuery combines user filter with performance-metrics constraint
func (s *PerformanceArtifactService) buildPerformanceFilterQuery(userFilter string) string {
	performanceFilter := "metricsType.string_value = 'performance-metrics'"
	if userFilter == "" {
		return performanceFilter
	}
	return fmt.Sprintf("(%s) AND (%s)", userFilter, performanceFilter)
}

func (s *PerformanceArtifactService) processArtifacts(artifacts []CatalogMetricsArtifact, params PerformanceArtifactParams) []CatalogMetricsArtifact {
	if params.Recommendations && params.TargetRPS <= 0 {
		params.TargetRPS = 1
	}

	// Set default property names if not provided
	rpsProperty := params.RPSProperty
	if rpsProperty == "" {
		rpsProperty = "requests_per_second"
	}
	latencyProperty := params.LatencyProperty
	if latencyProperty == "" {
		latencyProperty = "ttft_p90"
	}
	hardwareCountProperty := params.HardwareCountProperty
	if hardwareCountProperty == "" {
		hardwareCountProperty = "hardware_count"
	}
	hardwareTypeProperty := params.HardwareTypeProperty
	if hardwareTypeProperty == "" {
		hardwareTypeProperty = "hardware_type"
	}

	// Apply targetRPS calculations
	if params.TargetRPS > 0 {
		for i := range artifacts {
			s.addTargetRPSCalculations(artifacts[i], params.TargetRPS, rpsProperty)
		}
	}

	if params.Recommendations {
		artifacts = s.generateRecommendations(artifacts, latencyProperty, hardwareCountProperty, hardwareTypeProperty)
	}

	return artifacts
}

// addTargetRPSCalculations adds replicas and total_requests_per_second custom properties
func (s *PerformanceArtifactService) addTargetRPSCalculations(artifact CatalogMetricsArtifact, targetRPS int32, rpsProperty string) {
	customProperties := artifact.GetCustomProperties()
	if customProperties == nil {
		return
	}

	// Calculate replicas and total RPS based on performance metrics
	rps := s.extractCustomPropertiesDoubleValue(customProperties, rpsProperty, 1)
	if rps <= 0 {
		rps = 1
	}

	replicas := math.Ceil(float64(targetRPS) / rps)
	if replicas > math.MaxInt32 {
		replicas = math.MaxInt32
	} else if replicas < 1 {
		replicas = 1
	}
	totalRPS := replicas * rps

	// Initialize custom properties if nil
	if *customProperties == nil {
		*customProperties = []models.Properties{}
	}

	*customProperties = append(*customProperties,
		models.NewIntProperty("replicas", int32(replicas), true),
		models.NewDoubleProperty("total_requests_per_second", totalRPS, true))
}

const (
	initialEpsilon = 0.05
	fastEpsilon    = 0.1
	fastThreshold  = 100
)

// generateRecommendations removes duplicates based on cost estimates
func (s *PerformanceArtifactService) generateRecommendations(artifacts []CatalogMetricsArtifact, latencyProperty, hardwareCountProperty, hardwareTypeProperty string) []CatalogMetricsArtifact {
	keepIDs := map[int32]struct{}{}

	// Group the full list by hardware_type
	byHardware := s.groupArtifactsByStringProperty(artifacts, hardwareTypeProperty)

	for _, subArtifacts := range byHardware {
		// First pass, remove options that are slower and more expensive than
		// options seen before.
		s.sortArtifacts(subArtifacts, latencyProperty, hardwareCountProperty)
		cheapest := int32(math.MaxInt32)
		filtered := []CatalogMetricsArtifact{}
		for i, artifact := range subArtifacts {
			hwCount := s.extractCustomPropertiesIntValue(artifact.GetCustomProperties(), hardwareCountProperty, 1)
			replicas := s.extractCustomPropertiesIntValue(artifact.GetCustomProperties(), "replicas", 1)
			cost := hwCount * replicas

			// The list is sorted by latency, so every entry is at least as
			// slow as the previous entry and we therefore discard any
			// entry that isn't cheaper.
			if cost < cheapest {
				filtered = append(filtered, subArtifacts[i])
				cheapest = cost
			}
		}

		// The first pass gives us a list sorted by cost descending, reverse for cost ascending.
		slices.Reverse(filtered)

		// Second pass, keep options that offer a meaningful latency improvement over the last kept option.

		// epsilon is the threshold for an improvement to be considered meaningful.
		epsilon := initialEpsilon

		// Keep track of the IDs, because we want the output to be in the same order as the input to this function.
		id := filtered[0].GetID()
		if id == nil {
			return artifacts
		}
		keepIDs[*id] = struct{}{}

		lastKeptLatency := s.extractCustomPropertiesDoubleValue(filtered[0].GetCustomProperties(), latencyProperty, 0)

		for _, artifact := range filtered[1:] {
			if lastKeptLatency <= fastThreshold {
				// Increase the threshold once we see a "fast" option.
				epsilon = fastEpsilon
			}

			latency := s.extractCustomPropertiesDoubleValue(artifact.GetCustomProperties(), latencyProperty, lastKeptLatency)
			if (lastKeptLatency-latency)/lastKeptLatency > epsilon {
				id := artifact.GetID()
				if id != nil {
					keepIDs[*id] = struct{}{}
					lastKeptLatency = latency
				}
			}
		}
	}

	// Build the filtered list in the order we received it.
	result := make([]CatalogMetricsArtifact, 0, len(keepIDs))
	for _, artifact := range artifacts {
		id := artifact.GetID()
		if id == nil {
			continue
		}
		if _, ok := keepIDs[*id]; ok {
			result = append(result, artifact)
		}
	}
	return result
}

// validateCustomProperties ensures at least one artifact contains each custom property
// When there are no artifacts, validation passes (empty result is valid)
func (s *PerformanceArtifactService) validateCustomProperties(artifacts []CatalogMetricsArtifact, rpsProperty, latencyProperty, hardwareCountProperty, hardwareTypeProperty string) error {
	// If there are no artifacts, validation passes - empty result is valid
	if len(artifacts) == 0 {
		return nil
	}

	propertyExists := map[string]bool{
		rpsProperty:           false,
		latencyProperty:       false,
		hardwareCountProperty: false,
		hardwareTypeProperty:  false,
	}

	// Scan all artifacts to check property existence
	for _, artifact := range artifacts {
		props := artifact.GetCustomProperties()
		if props == nil {
			continue
		}

		for _, prop := range *props {
			if _, tracked := propertyExists[prop.Name]; tracked {
				propertyExists[prop.Name] = true
			}
		}
	}

	// Collect missing properties
	var missing []string
	for propName, exists := range propertyExists {
		if propName != "" && !exists {
			missing = append(missing, propName)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("properties not found in any performance artifacts: %v", missing)
	}

	return nil
}

func (s *PerformanceArtifactService) groupArtifactsByStringProperty(artifacts []CatalogMetricsArtifact, property string) map[string][]CatalogMetricsArtifact {
	results := map[string][]CatalogMetricsArtifact{}

	for i, artifact := range artifacts {
		value := s.extractCustomPropertiesStringValue(artifact.GetCustomProperties(), property)
		results[value] = append(results[value], artifacts[i])
	}

	return results
}

// sortArtifacts sorts a list of performance artifacts by latency, and hardware count as a tie-breaker.
func (s *PerformanceArtifactService) sortArtifacts(artifacts []CatalogMetricsArtifact, latencyProperty string, hardwareCountProperty string) {
	slices.SortFunc(artifacts, func(a CatalogMetricsArtifact, b CatalogMetricsArtifact) int {
		defaultLatency := math.MaxFloat64
		aLatency := s.extractCustomPropertiesDoubleValue(a.GetCustomProperties(), latencyProperty, defaultLatency)
		if aLatency <= 0 {
			aLatency = defaultLatency
		}
		bLatency := s.extractCustomPropertiesDoubleValue(b.GetCustomProperties(), latencyProperty, defaultLatency)
		if bLatency <= 0 {
			bLatency = defaultLatency
		}

		latencyCmp := cmp.Compare(aLatency, bLatency)
		if latencyCmp != 0 {
			return latencyCmp
		}

		// Against all odds, the latency values are the same, so prefer the one with a lower hardware count.
		var defaultCount int32 = math.MaxInt32
		aCount := s.extractCustomPropertiesIntValue(a.GetCustomProperties(), hardwareCountProperty, defaultCount)
		if aCount <= 0 {
			aCount = defaultCount
		}
		bCount := s.extractCustomPropertiesIntValue(b.GetCustomProperties(), hardwareCountProperty, defaultCount)
		if bCount <= 0 {
			bCount = defaultCount
		}

		return cmp.Compare(aCount, bCount)
	})
}

func (s *PerformanceArtifactService) extractCustomPropertiesDoubleValue(props *[]models.Properties, name string, def float64) float64 {
	value := def
	if props != nil {
		for _, prop := range *props {
			if prop.Name == name {
				if prop.DoubleValue != nil {
					value = *prop.DoubleValue
				}
				break
			}
		}
	}
	return value
}

func (s *PerformanceArtifactService) extractCustomPropertiesIntValue(props *[]models.Properties, name string, def int32) int32 {
	value := def
	if props != nil {
		for _, prop := range *props {
			if prop.Name == name {
				if prop.IntValue != nil {
					value = *prop.IntValue
				}
				break
			}
		}
	}
	return value
}

func (s *PerformanceArtifactService) extractCustomPropertiesStringValue(props *[]models.Properties, name string) string {
	if props != nil {
		for _, prop := range *props {
			if prop.Name == name {
				if prop.StringValue != nil {
					return *prop.StringValue
				}
				break
			}
		}
	}
	return ""
}

func (s *PerformanceArtifactService) paginate(list *models.ListWrapper[CatalogMetricsArtifact], pageSize int32, nextPageToken *string) {
	var cursor *scopes.Cursor
	if nextPageToken != nil {
		// If there's an error parsing the token, return the first page.
		cursor, _ = scopes.DecodeCursor(*nextPageToken)
	}

	if cursor != nil {
		var index int
		for index = range list.Items {
			id := list.Items[index].GetID()
			if id != nil && *id == cursor.ID {
				break
			}
		}
		list.Items = list.Items[index+1:]
	}

	list.NextPageToken = ""
	if len(list.Items) > int(pageSize) {
		list.Items = list.Items[:pageSize]

		lastID := list.Items[len(list.Items)-1].GetID()
		if lastID != nil {
			list.NextPageToken = scopes.CreateNextPageToken(*lastID, "")
		}
	}

	if len(list.Items) < math.MaxInt32 {
		list.Size = int32(len(list.Items))
	} else {
		list.Size = math.MaxInt32
	}
}
