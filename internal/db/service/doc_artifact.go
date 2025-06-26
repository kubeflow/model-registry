package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/scopes"
	"gorm.io/gorm"
)

var ErrDocArtifactNotFound = errors.New("doc artifact by id not found")

type DocArtifactRepositoryImpl struct {
	db     *gorm.DB
	typeID int64
}

func NewDocArtifactRepository(db *gorm.DB, typeID int64) models.DocArtifactRepository {
	return &DocArtifactRepositoryImpl{db: db, typeID: typeID}
}

func (r *DocArtifactRepositoryImpl) GetByID(id int32) (models.DocArtifact, error) {
	docArtifact := &schema.Artifact{}
	properties := []schema.ArtifactProperty{}

	if err := r.db.Where("id = ? AND type_id = ?", id, r.typeID).First(docArtifact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrDocArtifactNotFound, err)
		}

		return nil, fmt.Errorf("error getting doc artifact by id: %w", err)
	}

	if err := r.db.Where("artifact_id = ?", docArtifact.ID).Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("error getting properties by doc artifact id: %w", err)
	}

	return mapDataLayerToDocArtifact(*docArtifact, properties), nil
}

func (r *DocArtifactRepositoryImpl) Save(docArtifact models.DocArtifact, modelVersionID *int32) (models.DocArtifact, error) {
	now := time.Now().UnixMilli()

	docArtifactArt := mapDocArtifactToArtifact(docArtifact)
	properties := mapDocArtifactToArtifactProperties(docArtifact, docArtifactArt.ID)

	docArtifactArt.LastUpdateTimeSinceEpoch = now

	if docArtifact.GetID() == nil {
		glog.Info("Creating new DocArtifact")

		docArtifactArt.CreateTimeSinceEpoch = now
	} else {
		glog.Infof("Updating DocArtifact %d", *docArtifact.GetID())
	}

	hasCustomProperties := docArtifact.GetCustomProperties() != nil

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&docArtifactArt).Error; err != nil {
			return fmt.Errorf("error saving doc artifact: %w", err)
		}

		properties = mapDocArtifactToArtifactProperties(docArtifact, docArtifactArt.ID)
		existingCustomProperties := []schema.ArtifactProperty{}

		if err := tx.Where("artifact_id = ? AND is_custom_property = ?", docArtifactArt.ID, true).Find(&existingCustomProperties).Error; err != nil {
			return fmt.Errorf("error getting existing custom properties by doc artifact id: %w", err)
		}

		if hasCustomProperties {
			for _, existingProp := range existingCustomProperties {
				found := false
				for _, prop := range properties {
					if prop.Name == existingProp.Name && prop.ArtifactID == existingProp.ArtifactID && prop.IsCustomProperty == existingProp.IsCustomProperty {
						found = true
						break
					}
				}

				if !found {
					if err := tx.Delete(&existingProp).Error; err != nil {
						return fmt.Errorf("error deleting doc artifact property: %w", err)
					}
				}
			}
		}

		for _, prop := range properties {
			var existingProp schema.ArtifactProperty
			result := tx.Where("artifact_id = ? AND name = ? AND is_custom_property = ?",
				prop.ArtifactID, prop.Name, prop.IsCustomProperty).First(&existingProp)

			switch result.Error {
			case nil:
				prop.ArtifactID = existingProp.ArtifactID
				prop.Name = existingProp.Name
				prop.IsCustomProperty = existingProp.IsCustomProperty
				if err := tx.Model(&existingProp).Updates(prop).Error; err != nil {
					return fmt.Errorf("error updating doc artifact property: %w", err)
				}
			case gorm.ErrRecordNotFound:
				if err := tx.Create(&prop).Error; err != nil {
					return fmt.Errorf("error creating doc artifact property: %w", err)
				}
			default:
				return fmt.Errorf("error checking existing property: %w", result.Error)
			}
		}

		if modelVersionID != nil {
			// Check if attribution already exists to avoid duplicate key errors
			var existingAttribution schema.Attribution
			result := tx.Where("context_id = ? AND artifact_id = ?", *modelVersionID, docArtifactArt.ID).First(&existingAttribution)

			if result.Error == gorm.ErrRecordNotFound {
				// Attribution doesn't exist, create it
				attribution := schema.Attribution{
					ContextID:  *modelVersionID,
					ArtifactID: docArtifactArt.ID,
				}

				if err := tx.Create(&attribution).Error; err != nil {
					return fmt.Errorf("error creating attribution: %w", err)
				}
			} else if result.Error != nil {
				return fmt.Errorf("error checking existing attribution: %w", result.Error)
			}
			// If attribution already exists, do nothing
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mapDataLayerToDocArtifact(docArtifactArt, properties), nil
}

func (r *DocArtifactRepositoryImpl) List(listOptions models.DocArtifactListOptions) (*models.ListWrapper[models.DocArtifact], error) {
	list := models.ListWrapper[models.DocArtifact]{
		PageSize: listOptions.GetPageSize(),
	}

	docArtifacts := []models.DocArtifact{}
	docArtifactsArt := []schema.Artifact{}

	query := r.db.Model(&schema.Artifact{}).Where("type_id = ?", r.typeID)

	if listOptions.Name != nil {
		query = query.Where("name = ?", listOptions.Name)
	} else if listOptions.ExternalID != nil {
		query = query.Where("external_id = ?", listOptions.ExternalID)
	}

	if listOptions.ModelVersionID != nil {
		query = query.Joins("JOIN Attribution ON Attribution.artifact_id = Artifact.id").
			Where("Attribution.context_id = ?", listOptions.ModelVersionID)
		// Use table-prefixed pagination to avoid column ambiguity
		query = query.Scopes(scopes.PaginateWithTablePrefix(docArtifacts, &listOptions.Pagination, r.db, "Artifact"))
	} else {
		query = query.Scopes(scopes.Paginate(docArtifacts, &listOptions.Pagination, r.db))
	}

	if err := query.Find(&docArtifactsArt).Error; err != nil {
		return nil, fmt.Errorf("error listing doc artifacts: %w", err)
	}

	hasMore := false
	pageSize := listOptions.GetPageSize()
	if pageSize > 0 {
		hasMore = len(docArtifactsArt) > int(pageSize)
		if hasMore {
			docArtifactsArt = docArtifactsArt[:len(docArtifactsArt)-1]
		}
	}

	for _, docArtifactArt := range docArtifactsArt {
		properties := []schema.ArtifactProperty{}
		if err := r.db.Where("artifact_id = ?", docArtifactArt.ID).Find(&properties).Error; err != nil {
			return nil, fmt.Errorf("error getting properties by doc artifact id: %w", err)
		}

		docArtifact := mapDataLayerToDocArtifact(docArtifactArt, properties)
		docArtifacts = append(docArtifacts, docArtifact)
	}

	if hasMore && len(docArtifactsArt) > 0 {
		lastModel := docArtifactsArt[len(docArtifactsArt)-1]
		orderBy := listOptions.GetOrderBy()
		value := ""
		if orderBy != "" {
			switch orderBy {
			case "ID":
				value = fmt.Sprintf("%d", lastModel.ID)
			case "CREATE_TIME":
				value = fmt.Sprintf("%d", lastModel.CreateTimeSinceEpoch)
			case "LAST_UPDATE_TIME":
				value = fmt.Sprintf("%d", lastModel.LastUpdateTimeSinceEpoch)
			default:
				value = fmt.Sprintf("%d", lastModel.ID)
			}
		}
		nextToken := scopes.CreateNextPageToken(lastModel.ID, value)
		listOptions.NextPageToken = &nextToken
	} else {
		listOptions.NextPageToken = nil
	}

	list.Items = docArtifacts
	list.NextPageToken = listOptions.GetNextPageToken()
	list.PageSize = listOptions.GetPageSize()
	list.Size = int32(len(docArtifacts))

	return &list, nil
}

func mapDocArtifactToArtifact(docArtifact models.DocArtifact) schema.Artifact {
	if docArtifact == nil {
		return schema.Artifact{}
	}

	artifact := schema.Artifact{
		ID:     apiutils.ZeroIfNil(docArtifact.GetID()),
		TypeID: apiutils.ZeroIfNil(docArtifact.GetTypeID()),
	}

	if docArtifact.GetAttributes() != nil {
		artifact.Name = docArtifact.GetAttributes().Name
		artifact.URI = docArtifact.GetAttributes().URI
		artifact.ExternalID = docArtifact.GetAttributes().ExternalID
		if docArtifact.GetAttributes().State != nil {
			stateValue := models.Artifact_State_value[*docArtifact.GetAttributes().State]
			artifact.State = &stateValue
		}
		artifact.CreateTimeSinceEpoch = apiutils.ZeroIfNil(docArtifact.GetAttributes().CreateTimeSinceEpoch)
		artifact.LastUpdateTimeSinceEpoch = apiutils.ZeroIfNil(docArtifact.GetAttributes().LastUpdateTimeSinceEpoch)
	}

	return artifact
}

func mapDocArtifactToArtifactProperties(docArtifact models.DocArtifact, artifactID int32) []schema.ArtifactProperty {
	if docArtifact == nil {
		return []schema.ArtifactProperty{}
	}

	properties := []schema.ArtifactProperty{}

	if docArtifact.GetProperties() != nil {
		for _, prop := range *docArtifact.GetProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, false))
		}
	}

	if docArtifact.GetCustomProperties() != nil {
		for _, prop := range *docArtifact.GetCustomProperties() {
			properties = append(properties, mapPropertiesToArtifactProperty(prop, artifactID, true))
		}
	}

	return properties
}

func mapDataLayerToDocArtifact(docArtifact schema.Artifact, artProperties []schema.ArtifactProperty) models.DocArtifact {
	var state *string

	docArtifactType := models.DocArtifactType

	if docArtifact.State != nil {
		docState := models.Artifact_State_name[*docArtifact.State]
		state = &docState
	}

	docArtifactArt := models.DocArtifactImpl{
		ID:     &docArtifact.ID,
		TypeID: &docArtifact.TypeID,
		Attributes: &models.DocArtifactAttributes{
			Name:                     docArtifact.Name,
			URI:                      docArtifact.URI,
			State:                    state,
			ArtifactType:             &docArtifactType,
			ExternalID:               docArtifact.ExternalID,
			CreateTimeSinceEpoch:     &docArtifact.CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: &docArtifact.LastUpdateTimeSinceEpoch,
		},
	}

	customProperties := []models.Properties{}
	properties := []models.Properties{}

	for _, prop := range artProperties {
		if prop.IsCustomProperty {
			customProperties = append(customProperties, mapDataLayerToArtifactProperties(prop))
		} else {
			properties = append(properties, mapDataLayerToArtifactProperties(prop))
		}
	}

	docArtifactArt.CustomProperties = &customProperties
	docArtifactArt.Properties = &properties

	return &docArtifactArt
}
