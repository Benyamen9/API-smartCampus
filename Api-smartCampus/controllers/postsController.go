package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	models "github.com/jalil/Api-smartCampus/Models"
	"github.com/jalil/Api-smartCampus/initializers"
)

// Définition des structures GeoJSON
type Geometry struct {
	Type        string   `json:"type"`
	Coordinates []string `json:"coordinates"`
}

type Properties struct {
	Time   string  `json:"time"`
	Value  float64 `json:"value"`
	Symbol string  `json:"symbol"`
}

type History struct {
	Time   string  `json:"time"`
	Value  float64 `json:"value"`
	Symbol string  `json:"symbol"`
}

type Feature struct {
	Source      string     `json:"source"`
	SourceID    int        `json:"sourceId"`
	Geometry    Geometry   `json:"geometry"`
	Properties  Properties `json:"properties"`
	History     []History  `json:"history"`
	LastUpdated string     `json:"lastUpdated"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// Structure du fichier JSON energy.json
type EnergyData struct {
	Results []struct {
		Series []struct {
			Columns []string        `json:"columns"`
			Values  [][]interface{} `json:"values"`
		} `json:"series"`
	} `json:"results"`
}

func TabsensorIndex(c *gin.Context) {

	// Récupération des données de la base PostgreSQL
	var tabsensors []models.Tabsensor
	result := initializers.DB.Find(&tabsensors)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des données depuis PostgreSQL"})
		return
	}

	// Lecture du fichier JSON
	filePath := "data/energy.json"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erreur d'ouverture du fichier JSON:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Fichier JSON introuvable"})
		return
	}
	defer file.Close()

	// Décoder le JSON
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de lire le fichier JSON"})
		return
	}

	var energyData EnergyData
	err = json.Unmarshal(data, &energyData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Format JSON invalide"})
		return
	}

	// Fusion des données PostgreSQL et JSON
	var features []Feature

	for _, sensor := range tabsensors {
		var history []History

		// Recherche des valeurs correspondantes dans le JSON
		for _, result := range energyData.Results {
			for _, series := range result.Series {
				for _, value := range series.Values {
					if len(value) >= 5 {
						//fmt.Println("🔍 Valeur trouvée :", value)

						// ✅ Vérifie que `sourceid` est bien un string et convertit en int
						sourceIDStr, ok := value[2].(string)
						if !ok {
							fmt.Println("❌ Erreur: sourceID n'est pas un string:", value[2])
							continue
						}

						sourceID, err := strconv.Atoi(sourceIDStr)
						if err != nil {
							fmt.Println("❌ Erreur conversion sourceID:", err)
							continue
						}

						// ✅ Vérifie si le sourceID correspond au capteur
						if sourceID == sensor.Sourceid {

							// ✅ Vérifie que `value[4]` est bien un float64
							var valueFloat float64
							switch v := value[4].(type) {
							case float64:
								valueFloat = v
							case string:
								valueFloat, err = strconv.ParseFloat(v, 64)
								if err != nil {
									fmt.Println("❌ Erreur conversion value[4]:", err)
									continue
								}
							default:
								fmt.Println("❌ Type inconnu pour value[4]:", value[4])
								continue
							}

							history = append(history, History{
								Time:   value[0].(string),
								Value:  valueFloat,
								Symbol: value[3].(string),
							})
						}
					}
				}
			}
		}

		// Si l'historique est vide, on met des valeurs par défaut
		var properties Properties
		if len(history) > 0 {
			properties = Properties{
				Time:   history[0].Time,
				Value:  history[0].Value,
				Symbol: history[0].Symbol,
			}
		} else {
			properties = Properties{
				Time:   "N/A",
				Value:  0.0,
				Symbol: "N/A",
			}
		}

		// Création de l'objet Feature
		feature := Feature{
			Source:   "PRODUCTION PHOTOVOLTAÏQUE",
			SourceID: sensor.Sourceid,
			Geometry: Geometry{
				Type: "Point",
				Coordinates: []string{
					strings.TrimSpace(sensor.Latitude),
					strings.TrimSpace(sensor.Longitude),
				},
			},
			Properties:  properties,
			History:     history,
			LastUpdated: "N/A",
		}

		features = append(features, feature)
	}

	// Construction de la réponse GeoJSON
	response := FeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Envoi de la réponse
	c.JSON(http.StatusOK, response)
}

func JsonDataGet(c *gin.Context) {
	filePath := "data/energy.json"

	if _, err := os.Stat(filePath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fichier introuvable"})
		return
	}

	c.File(filePath)
}

func TabsensorShow(c *gin.Context) {
	// Get id from url

	id := c.Param("id")

	// Get tabsensor
	var tabsensor models.Tabsensor
	initializers.DB.First(&tabsensor, id)

	// Respond with them
	c.JSON(200, gin.H{
		"tabsensor": tabsensor,
	})

}

/*
func TabsensorUpdate(c *gin.Context) {
	// Get the id
	id := c.Param("id")

	//Get the data
	var body struct {
		SourceID  int
		Latitude  string
		Longitude string
	}

	c.Bind(&body)

	// Find the post
	var tabsensor models.Tabsensor
	initializers.DB.First(&tabsensor, id)

	// Update it
	initializers.DB.Model(&tabsensor).Updates(models.Tabsensor{
		Sourceid:  body.SourceID,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	})

	// Respond with it
	c.JSON(200, gin.H{
		"tabsensor": tabsensor,
	})

}
*/

/*
func TabsensorDelete(c *gin.Context) {

	// Get the id
	id := c.Param("id")

	// Delete the post
	initializers.DB.Delete(&models.Tabsensor{}, id)

	// Respond with 200 OK
	c.Status(200)

}
*/
