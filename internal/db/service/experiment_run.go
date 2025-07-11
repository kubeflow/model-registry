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

var ErrExperimentRunNotFound = errors.New("experiment run by id not found")

type ExperimentRunRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewExperimentRunRepository(db *gorm.DB, typeID int64) models.ExperimentRunRepository {
	return &ExperimentRunRepositoryImpl{db: db, typeID: typeID}
}

func (r *ExperimentRunRepositoryImpl) GetByID(id int32) (models.ExperimentRun, error) {
	expRunCtx := &schema.Context{}
	propertiesCtx := []schema.ContextProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(expRunCtx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrExperimentRunNotFound, err)
		}

		return nil, fmt.Errorf("error getting experiment run by id: %w", err)
	}

	if err := r.db.Where("context_id = ?", expRunCtx.ID).Find(&propertiesCtx).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by experiment run id: %w", err)
	}

	return mapDataLayerToExperimentRun(*expRunCtx, propertiesCtx), nil
}

func (r *ExperimentRunRepositoryImpl) Save(experimentRun models.ExperimentRun, experimentID *int32) (models.ExperimentRun, error) {
	now := time.Now().UnixMilli()

	expRunCtx := mapExperimentRunToContext(experimentRun)
	var finalPropertiesCtx []schema.ContextProperty

	expRunCtx.LastUpdateTimeSinceEpoch = now

	if experimentRun.GetID() == nil {
		expRunCtx.CreateTimeSinceEpoch = now
	}

	hasCustomProperties := experimentRun.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&expRunCtx).Error; err != nil {
			return fmt.Errorf("error saving experiment run context: %w", err)
		}

		propertiesCtx := mapExperimentRunToContextProperties(experimentRun, expRunCtx.ID)
		existingCustomPropertiesCtx := []schema.ContextProperty{}

		if err := tx.Where("context_id = ? AND is_custom_property = ?", expRunCtx.ID, true).Find(&existingCustomPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by experiment run id: %w", err)
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
						return fmt.Errorf("error deleting experiment run context property: %w", err)
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
					return fmt.Errorf("error updating experiment run context property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating experiment run context property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		// Handle parent-child relationship with experiment
		if experimentID != nil {
			parentsCtx := schema.ParentContext{
				ContextID:       expRunCtx.ID,
				ParentContextID: *experimentID,
			}

			if err := tx.Save(&parentsCtx).Error; err != nil {
				return fmt.Errorf("error saving experiment run parent context: %w", err)
			}
		}

		// Get all final properties for the return object
		if err := tx.Where("context_id = ?", expRunCtx.ID).Find(&finalPropertiesCtx).Error; err != nil {
			return fmt.Errorf("error getting final properties by experiment run id: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Return the updated entity using the data we already have
	return mapDataLayerToExperimentRun(expRunCtx, finalPropertiesCtx), nil
}

func (r *ExperimentRunRepositoryImpl) List(listOptions models.ExperimentRunListOptions) (*models.ListWrapper[models.ExperimentRun], error) {
	list := models.ListWrapper[models.ExperimentRun]{
		PageSize: listOptions.GetPageSize(),
	}

	experimentRuns := []models.ExperimentRun{}
	experimentRunsCtx := []schema.Context{}

	query := r.db.Model(&schema.Context{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ExperimentID != nil {
		query = query.Joins("JOIN ParentContext ON ParentContext.context_id = Context.id").
			Where("ParentContext.parent_context_id = ?", listOptions.ExperimentID)
		query = query.Scopes(scopes.PaginateWithTablePrefix(experimentRuns, &listOptions.Pagination, r.db, "Context"))
	} else {
		query = query.Scopes(scopes.Paginate(experimentRuns, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&experimentRunsCtx).Error; err != nil {
		return nil, fmt.Errorf("error listing experiment runs: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(experimentRunsCtx) > int(pageSize)
		if hasMore {
			experimentRunsCtx = experimentRunsCtx[:len(experimentRunsCtx)-1]
		}
	}

	for _, expRunCtx := range experimentRunsCtx {
		propertiesCtx := []schema.ContextProperty{}
		if err := r.db.Where("context_id = ?", expRunCtx.ID).Find(&propertiesCtx).Error; err != nil {
			return nil, fmt.Errorf("error getting properties for experiment run %d: %w", expRunCtx.ID, err)
		}
		experimentRun := mapDataLayerToExperimentRun(expRunCtx, propertiesCtx)
		experimentRuns = append(experimentRuns, experimentRun)
	}

	if hasMore && len(experimentRunsCtx) > 0 {
		lastExperimentRun := experimentRunsCtx[len(experimentRunsCtx)-1]
		orderBy := listOptions.GetOrderBy()
		value := ""
		if orderBy != "" {
			switch orderBy {
			case "ID":
				value = fmt.Sprintf("%d", lastExperimentRun.ID)
			case "CREATE_TIME":
				value = fmt.Sprintf("%d", lastExperimentRun.CreateTimeSinceEpoch)
			case "LAST_UPDATE_TIME":
				value = fmt.Sprintf("%d", lastExperimentRun.LastUpdateTimeSinceEpoch)
			default:
				value = fmt.Sprintf("%d", lastExperimentRun.ID)
			}
		}
		nextToken := scopes.CreateNextPageToken(lastExperimentRun.ID, value)
		listOptions.NextPageToken = &nextToken
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = experimentRuns
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(experimentRuns))

	return &list, nil
}

func mapExperimentRunToContext(experimentRun models.ExperimentRun) schema.Context {
	if experimentRun == nil {
		return schema.Context{}
	}

	expRunCtx := schema.Context{
		ID:     apiutils.ZeroIfNil(experimentRun.GetID()),
		TypeID: apiutils.ZeroIfNil(experimentRun.GetTypeID()),
	}

	if experimentRun.GetAttributes() != nil {
		expRunCtx.Name = apiutils.ZeroIfNil(experimentRun.GetAttributes().Name)
		expRunCtx.ExternalID = experimentRun.GetAttributes().ExternalID
		expRunCtx.CreateTimeSinceEpoch = apiutils.ZeroIfNil(experimentRun.GetAttributes().CreateTimeSinceEpoch)
		expRunCtx.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(experimentRun.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return expRunCtx
}

func mapExperimentRunToContextProperties(experimentRun models.ExperimentRun, experimentRunId int32) []schema.ContextProperty {
	if experimentRun == nil {
		return []schema.ContextProperty{}
	}

	propertiesCtx := []schema.ContextProperty{}

	if experimentRun.GetProperties() != nil {
		for _, prop := range *experimentRun.GetProperties() {
			propCtx := mapPropertiesToContextProperty(prop, experimentRunId, false)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	if experimentRun.GetCustomProperties() != nil {
		for _, prop := range *experimentRun.GetCustomProperties() {
			propCtx := mapPropertiesToContextProperty(prop, experimentRunId, true)
			propertiesCtx = append(propertiesCtx, propCtx)
		}
	}

	return propertiesCtx
}

func mapDataLayerToExperimentRun(expRunCtx schema.Context, propertiesCtx []schema.ContextProperty) models.ExperimentRun {
	experimentRun := models.ExperimentRunImpl{
		ID:     &expRunCtx.ID,
		TypeID: &expRunCtx.TypeID,
		Attributes: &models.ExperimentRunAttributes{
			Name:                     &expRunCtx.Name,
			ExternalID:               expRunCtx.ExternalID,
			CreateTimeSinceEpoch:     &expRunCtx.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &expRunCtx.LastUpdateTimeSinceEpoch,
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

	experimentRun.CustomProperties = &customProperties
	experimentRun.Properties = &properties

	return &experimentRun
}
