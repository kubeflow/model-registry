package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"gorm.io/gorm"
)

var ErrExperimentNotFound = errors.New("experiment by id not found")

type ExperimentRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewExperimentRepository(db *gorm.DB, typeID int64) models.ExperimentRepository {
	return &ExperimentRepositoryImpl{db: db, typeID: typeID}
}

func (r *ExperimentRepositoryImpl) GetByID(id int32) (models.Experiment, error) {
	expCtx := &schema.Context{}
	propertiesCtx := []schema.ContextProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(expCtx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrExperimentNotFound, err)
		}

		return nil, fmt.Errorf("error getting experiment by id: %w", err)
	}

	if err := r.db.Where("context_id = ?", expCtx.ID).Find(&propertiesCtx).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by experiment id: %w", err)
	}

	return mapDataLayerToExperiment(*expCtx, propertiesCtx), nil
}

func (r *ExperimentRepositoryImpl) Save(experiment models.Experiment) (models.Experiment, error) {
	now := time.Now().UnixMilli()

	expCtx := mapExperimentToContext(experiment)
	var finalPropertiesCtx []schema.ContextProperty

	expCtx.LastUpdateTimeSinceEpoch = now

	if experiment.GetID() == nil {
		expCtx.CreateTimeSinceEpoch = now
	}

	hasCustomProperties := experiment.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&expCtx).Error; err != nil {
			return fmt.Errorf("error saving experiment context: %w", err)
		}

		propertiesCtx := mapExperimentToContextProperties(experiment, expCtx.ID)
		existingCustomPropertiesCtx := []schema.ContextProperty{}

		if err := tx.Where("context_id = ? AND is_custom_property = ?", expCtx.ID, true).Find(&existingCustomPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by experiment id: %w", err)
		}

		if hasCustomProperties {
			for _, existingProp := range existingCustomPropertiesCtx {
				found := false
				for _, prop := range propertiesCtx {
					if prop.Name == existingProp.Name && prop.ContextID == existingProp.ContextID && prop.IsCustomProperty == existingProp.IsCustomProperty {
						found = true
						break
					}
				}

				if !found {
					if err := tx.Delete(&existingProp).Error; err != nil {
						return fmt.Errorf("error deleting experiment context property: %w", err)
					}
				}
			}
		}

		for _, prop := range propertiesCtx {
			var existingProp schema.ContextProperty
			result := tx.Where("context_id = ? AND name = ? AND is_custom_property = ?",
				prop.ContextID, prop.Name, prop.IsCustomProperty).First(&existingProp)

			switch result.Error {
			case nil:
				prop.ContextID = existingProp.ContextID
				prop.Name = existingProp.Name
				prop.IsCustomProperty = existingProp.IsCustomProperty
				if err := tx.Model(&existingProp).Updates(prop).Error; err != nil {
					return fmt.Errorf("error updating experiment context property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating experiment context property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		// Get all final properties for the return object
		if err := tx.Where("context_id = ?", expCtx.ID).Find(&finalPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting final properties by experiment id: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Return the updated entity using the data we already have
	return mapDataLayerToExperiment(expCtx, finalPropertiesCtx), nil
}

func (r *ExperimentRepositoryImpl) List(listOptions models.ExperimentListOptions) (*models.ListWrapper[models.Experiment], error) {
	list := models.ListWrapper[models.Experiment]{
		PageSize: listOptions.GetPageSize(),
	}

	experiments := []models.Experiment{}
	experimentsCtx := []schema.Context{}

	query := r.db.Model(&schema.Context{}).Where("type_id = ?", r.typeID)
	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	query = query.Scopes(scopes.Paginate(experiments, &listOptions.Pagination, r.db))

	if err := query.Find(&experimentsCtx).Error; err != nil {
		return nil, fmt.Errorf("error listing experiments: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(experimentsCtx) > int(pageSize)
		if hasMore {
			experimentsCtx = experimentsCtx[:len(experimentsCtx)-1]
		}
	}

	for _, expCtx := range experimentsCtx {
		propertiesCtx := []schema.ContextProperty{}
		if err := r.db.Where("context_id = ?", expCtx.ID).Find(&propertiesCtx).Error; err != nil {
			return nil, fmt.Errorf("error getting properties for experiment %d: %w", expCtx.ID, err)
		}
		experiment := mapDataLayerToExperiment(expCtx, propertiesCtx)
		experiments = append(experiments, experiment)
	}

	if hasMore && len(experimentsCtx) > 0 {
		lastExperiment := experimentsCtx[len(experimentsCtx)-1]
		orderBy := listOptions.GetOrderBy()
		value := ""
		if orderBy != "" {
			switch orderBy {
			case "ID":
				value = fmt.Sprintf("%d", lastExperiment.ID)
			case "CREATE_TIME":
				value = fmt.Sprintf("%d", lastExperiment.CreateTimeSinceEpoch)
			case "LAST_UPDATE_TIME":
				value = fmt.Sprintf("%d", lastExperiment.LastUpdateTimeSinceEpoch)
			default:
				value = fmt.Sprintf("%d", lastExperiment.ID)
			}
		}
		nextToken := scopes.CreateNextPageToken(lastExperiment.ID, value)
		listOptions.NextPageToken = &nextToken
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = experiments
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(experiments))

	return &list, nil
}

func mapExperimentToContext(experiment models.Experiment) schema.Context {
	if experiment == nil {
		return schema.Context{}
	}

	expCtx := schema.Context{
		ID:     apiutils.ZeroIfNil(experiment.GetID()),
		TypeID: apiutils.ZeroIfNil(experiment.GetTypeID()),
	}

	if experiment.GetAttributes() != nil {
		expCtx.Name = apiutils.ZeroIfNil(experiment.GetAttributes().Name)
		expCtx.ExternalID = experiment.GetAttributes().ExternalID
		expCtx.CreateTimeSinceEpoch = apiutils.ZeroIfNil(experiment.GetAttributes().CreateTimeSinceEpoch)
		expCtx.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(experiment.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return expCtx
}

func mapExperimentToContextProperties(experiment models.Experiment, experimentId int32) []schema.ContextProperty {
	if experiment == nil {
		return []schema.ContextProperty{}
	}

	propertiesCtx := []schema.ContextProperty{}

	if experiment.GetProperties() != nil {
		for _, prop := range *experiment.GetProperties() {
			propCtx := mapPropertiesToContextProperty(prop, experimentId, false)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	if experiment.GetCustomProperties() != nil {
		for _, prop := range *experiment.GetCustomProperties() {
			propCtx := mapPropertiesToContextProperty(prop, experimentId, true)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	return propertiesCtx
}

func mapDataLayerToExperiment(expCtx schema.Context, propertiesCtx []schema.ContextProperty) models.Experiment {
	experiment := models.ExperimentImpl{
		ID:     &expCtx.ID,
		TypeID: &expCtx.TypeID,
		Attributes: &models.ExperimentAttributes{
			Name:                     &expCtx.Name,
			ExternalID:               expCtx.ExternalID,
			CreateTimeSinceEpoch:     &expCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &expCtx.LastUpdateTimeSinceEpoch,
		},
	}

	customProperties := []models.Properties{}
	properties := []models.Properties{}

	for _, prop := range propertiesCtx {
		if prop.IsCustomProperty {
			customProperties = append(customProperties, mapDataLayerToProperties(prop))
		} else {
			properties = append(properties, mapDataLayerToProperties(prop))
		}
	}

	experiment.CustomProperties = &customProperties
	experiment.Properties = &properties

	return &experiment
}
