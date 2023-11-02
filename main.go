package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "postgres"
	dbname   = "OfficeManagement"
)

type OfficeMaster struct {
	OfficeID          int       `json:"OfficeID"`
	OfficeTypeID      int       `json:"OfficeTypeID"`
	OfficeName        string    `json:"OfficeName"`
	EmailID           string    `json:"EmailID"`
	ContactNumber     string    `json:"ContactNumber"`
	WorkingHoursFrom  time.Time `json:"WorkingHoursFrom"`
	WorkingHoursTo    time.Time `json:"WorkingHoursTo"`
	DivisionID        int       `json:"DivisionID"`
	RegionID          int       `json:"RegionID"`
	CircleID          int       `json:"CircleID"`
	ReportingOfficeID int64     `json:"ReportingOfficeID"`
	Latitude          float64   `json:"Latitude"`
	Longitude         float64   `json:"Longitude"`
	Status            string    `json:"Status"`
	CSIFacilityID     string    `json:"CSIFacilityID"`
	OpenToPublicDate  string    `json:"OpenToPublicDate"`
	ClosedDate        string    `json:"ClosedDate"`
	ReasonForDisable  string    `json:"ReasonForDisable"`
	ReasonToEnable    string    `json:"ReasonToEnable"`
	CreatedBy         string    `json:"CreatedBy"`
	CreatedDate       time.Time `json:"CreatedDate"`
	UpdatedBy         string    `json:"UpdatedBy"`
	UpdatedDate       time.Time `json:"UpdatedDate"`
	ValidatedFlag     string    `json:"ValidatedFlag"`
}

type OfficeAttributeData struct {
	AttributeID            int       `json:"AttributeID"`
	OfficeID               int       `json:"OfficeID"`
	OfficeTypeID           int       `json:"OfficeTypeID"`
	OpenedDate             time.Time `json:"OpenedDate"`
	ClosedDate             time.Time `json:"ClosedDate"`
	QRTerminalID           string    `json:"QRTerminalID"`
	OfficeAddressLine1     string    `json:"OfficeAddressLine1"`
	OfficeAddressLine2     string    `json:"OfficeAddressLine2"`
	OfficeAddressLine3     string    `json:"OfficeAddressLine3"`
	Landmark               string    `json:"Landmark"`
	CityID                 int       `json:"CityID"`
	DistrictID             int       `json:"DistrictID"`
	TalukID                int       `json:"TalukID"`
	VillageID              int       `json:"VillageID"`
	StateID                int       `json:"StateID"`
	Pincode                string    `json:"Pincode"`
	PAOCode                string    `json:"PAOCode"`
	SolId                  string    `json:"SolId"`
	PLIId                  string    `json:"PLIId"`
	GSTNForHO              string    `json:"GSTNForHO"`
	WEGCode                string    `json:"WEGCode"`
	DDOCode                string    `json:"DDOCode"`
	DeliveryOfficeFlag     bool      `json:"DeliveryOfficeFlag"`
	CSIRolledOutFlag       bool      `json:"CSIRolledOutFlag"`
	SingleHandedOfficeFlag bool      `json:"SingleHandedOfficeFlag"`
	CreatedBy              string    `json:"CreatedBy"`
	CreatedDate            time.Time `json:"CreatedDate"`
	UpdatedBy              string    `json:"UpdatedBy"`
	UpdatedDate            time.Time `json:"UpdatedDate"`
}

func logRequest(handler gin.HandlerFunc, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request details
		logger.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL)

		// Capture the start time
		startTime := time.Now()

		// Call the actual handler
		handler(c)

		// Calculate and log the request processing time
		elapsed := time.Since(startTime)
		logger.Printf("Request processed in %s\n", elapsed)
	}
}

func main() {
	// Create a log file
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer logFile.Close()

	// Initialize the logger to write to the log file
	logger := log.New(logFile, "", log.LstdFlags)

	// Connect to the database
	postgresqldbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := sql.Open("postgres", postgresqldbInfo)
	if err != nil {
		logger.Fatal("Database connection error:", err)
	}
	defer db.Close()
	logger.Printf("Database Connection successfully established %v\n", dbname)

	// Create a new Gin router
	r := gin.Default()

	// Define routes with logging
	r.GET("/officetypes", logRequest(getOfficeTypesHandler(db), logger))
	r.GET("/circles", logRequest(getCircleNameHandler(db), logger))
	r.GET("/regions", logRequest(getRegionsForCircleHandler(db), logger))
	r.GET("/divisions", logRequest(getDivisionsForRegionHandler(db), logger))
	r.GET("/subdivisions", logRequest(getSubDivisionsForDivisionHandler(db), logger))
	r.POST("/createoffice", logRequest(createOfficeHandler(db), logger))
	r.POST("/createofficeattributes", logRequest(createOfficeAttributeHandler(db), logger))
	r.PUT("/updateoffice/:OfficeID", logRequest(updateOfficeHandler(db), logger))
	r.PUT("/updateofficeattribute/:AttributeID", logRequest(updateOfficeAttributeHandler(db), logger))

	// Start the server
	port := 5032
	logger.Printf("Server started on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Fatal("Server error:", err)
	}
}

func getOfficeTypesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use Squirrel to build the query
		query := squirrel.Select("OfficeTypeCode", "OfficeTypeDescription").From("OfficeTypeMaster")

		// Get the SQL query and arguments
		sql, args, err := query.ToSql()
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Execute the query
		rows, err := db.Query(sql, args...)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var officeTypeData []map[string]interface{}

		// Iterate through the rows and build a response
		for rows.Next() {
			var officeTypeCode string
			var officeTypeDescription string
			if err := rows.Scan(&officeTypeCode, &officeTypeDescription); err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			officeType := map[string]interface{}{"office_type_code": officeTypeCode, "office_type_description": officeTypeDescription}
			officeTypeData = append(officeTypeData, officeType)
		}

		// Check for errors during iteration
		if err := rows.Err(); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Return the office type data as a JSON response
		c.JSON(http.StatusOK, officeTypeData)
	}
}
func getCircleNameHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute a query to retrieve CircleID and CircleName from the CircleMaster table
		rows, err := db.Query("SELECT CircleID, CircleName FROM CircleMaster")
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var circleData []map[string]interface{}

		// Iterate through the rows and build a response
		for rows.Next() {
			var circleID int
			var circleName string
			if err := rows.Scan(&circleID, &circleName); err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			circle := map[string]interface{}{"circle_id": circleID, "circle_name": circleName}
			circleData = append(circleData, circle)
		}

		// Return the circle data as a JSON response
		c.JSON(http.StatusOK, circleData)
	}
}
func getRegionsForCircleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the CircleName from the query parameters
		circleName := c.DefaultQuery("circleName", "")

		// Execute a query to retrieve region IDs and names based on CircleName
		rows, err := db.Query("SELECT RegionID, RegionName FROM RegionMaster WHERE CircleID = (SELECT CircleID FROM CircleMaster WHERE CircleName = $1)", circleName)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var regionData []map[string]interface{}

		// Iterate through the rows and build a response
		for rows.Next() {
			var regionID int
			var regionName string
			if err := rows.Scan(&regionID, &regionName); err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// Create a map for each region with RegionID and RegionName
			regionMap := map[string]interface{}{
				"region_id":   regionID,
				"region_name": regionName,
			}

			regionData = append(regionData, regionMap)
		}

		// Return the region data as a JSON response
		c.JSON(http.StatusOK, regionData)
	}
}

func getDivisionsForRegionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the RegionName from the query parameters
		regionName := c.DefaultQuery("regionName", "")

		// Execute a query to retrieve division IDs and names based on RegionName
		rows, err := db.Query("SELECT DivisionID, DivisionName FROM DivisionMaster WHERE RegionID = (SELECT RegionID FROM RegionMaster WHERE RegionName = $1)", regionName)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var divisionData []map[string]interface{}

		// Iterate through the rows and build a response
		for rows.Next() {
			var divisionID int
			var divisionName string
			if err := rows.Scan(&divisionID, &divisionName); err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// Create a map for each division with DivisionID and DivisionName
			divisionMap := map[string]interface{}{
				"division_id":   divisionID,
				"division_name": divisionName,
			}

			divisionData = append(divisionData, divisionMap)
		}

		// Return the division data as a JSON response
		c.JSON(http.StatusOK, divisionData)
	}
}

func getSubDivisionsForDivisionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the DivisionName from the query parameters
		divisionName := c.DefaultQuery("divisionName", "")

		// Execute a query to retrieve SubDivision IDs and names based on DivisionName
		rows, err := db.Query("SELECT SubDivisionID, SubDivisionName FROM SubDivisionMaster WHERE DivisionID = (SELECT DivisionID FROM DivisionMaster WHERE DivisionName = $1)", divisionName)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var subdivisionData []map[string]interface{}

		// Iterate through the rows and build a response
		for rows.Next() {
			var subDivisionID int
			var subDivisionName string
			if err := rows.Scan(&subDivisionID, &subDivisionName); err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// Create a map for each subdivision with SubDivisionID and SubDivisionName
			subdivisionMap := map[string]interface{}{
				"subdivision_id":   subDivisionID,
				"subdivision_name": subDivisionName,
			}

			subdivisionData = append(subdivisionData, subdivisionMap)
		}

		// Return the subdivision data as a JSON response
		c.JSON(http.StatusOK, subdivisionData)
	}
}
func createOfficeHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var officeData OfficeMaster
		if err := c.ShouldBindJSON(&officeData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON request"})
			return
		}

		_, err := db.Exec(`
			INSERT INTO OfficeMaster (
				OfficeTypeID, OfficeName, EmailID, ContactNumber, WorkingHoursFrom, WorkingHoursTo, DivisionID, RegionID, CircleID, ReportingOfficeId, Latitude, Longitude, Status, CSIFacilityID, OpenToPublicDate, ClosedDate, ReasonForDisable, ReasonToEnable, CreatedBy, CreatedDate, UpdatedBy, UpdatedDate, ValidatedFlag
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)`,
			officeData.OfficeTypeID, officeData.OfficeName, officeData.EmailID, officeData.ContactNumber, officeData.WorkingHoursFrom, officeData.WorkingHoursTo, officeData.DivisionID, officeData.RegionID, officeData.CircleID, officeData.ReportingOfficeID, officeData.Latitude, officeData.Longitude, officeData.Status, officeData.CSIFacilityID, officeData.OpenToPublicDate, officeData.ClosedDate, officeData.ReasonForDisable, officeData.ReasonToEnable, officeData.CreatedBy, officeData.CreatedDate, officeData.UpdatedBy, officeData.UpdatedDate, officeData.ValidatedFlag)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into the database"})
			return
		}

		// Send a success response
		c.JSON(http.StatusCreated, gin.H{"message": "Office created successfully"})
	}
}

func createOfficeAttributeHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var officeAttributeData OfficeAttributeData
		if err := c.ShouldBindJSON(&officeAttributeData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON request"})
			return
		}

		_, err := db.Exec(`
			INSERT INTO OfficeAttributeMaster (
				OfficeID, OfficeTypeID, OpenedDate, ClosedDate, QRTerminalID, OfficeAddressLine1, OfficeAddressLine2,
				OfficeAddressLine3, Landmark, CityID, DistrictID, TalukID, VillageID, StateID, Pincode, PAOCode, SolId,
				PLIId, GSTNForHO, WEGCode, DDOCode, DeliveryOfficeFlag, CSIRolledOutFlag, SingleHandedOfficeFlag,
				CreatedBy, CreatedDate, UpdatedBy, UpdatedDate
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,
				$23, $24, $25, $26, $27, $28
			)`,
			officeAttributeData.OfficeID, officeAttributeData.OfficeTypeID, officeAttributeData.OpenedDate,
			officeAttributeData.ClosedDate, officeAttributeData.QRTerminalID, officeAttributeData.OfficeAddressLine1,
			officeAttributeData.OfficeAddressLine2, officeAttributeData.OfficeAddressLine3, officeAttributeData.Landmark,
			officeAttributeData.CityID, officeAttributeData.DistrictID, officeAttributeData.TalukID, officeAttributeData.VillageID,
			officeAttributeData.StateID, officeAttributeData.Pincode, officeAttributeData.PAOCode, officeAttributeData.SolId,
			officeAttributeData.PLIId, officeAttributeData.GSTNForHO, officeAttributeData.WEGCode, officeAttributeData.DDOCode,
			officeAttributeData.DeliveryOfficeFlag, officeAttributeData.CSIRolledOutFlag, officeAttributeData.SingleHandedOfficeFlag,
			officeAttributeData.CreatedBy, officeAttributeData.CreatedDate, officeAttributeData.UpdatedBy, officeAttributeData.UpdatedDate)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into the database"})
			return
		}

		// Send a success response
		c.JSON(http.StatusCreated, gin.H{"message": "Office attribute created successfully"})
	}
}
func updateOfficeHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		officeID := c.Param("OfficeID")

		var officeData OfficeMaster
		if err := c.ShouldBindJSON(&officeData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON request"})
			return
		}

		_, err := db.Exec(`
            UPDATE OfficeMaster
            SET OfficeTypeID = $1, OfficeName = $2, EmailID = $3, ContactNumber = $4, WorkingHoursFrom = $5,
                WorkingHoursTo = $6, DivisionID = $7, RegionID = $8, CircleID = $9, ReportingOfficeId = $10,
                Latitude = $11, Longitude = $12, Status = $13, CSIFacilityID = $14, OpenToPublicDate = $15,
                ClosedDate = $16, ReasonForDisable = $17, ReasonToEnable = $18, CreatedBy = $19, UpdatedBy = $20,
                UpdatedDate = $21, ValidatedFlag = $22
            WHERE OfficeID = $23`,
			officeData.OfficeTypeID, officeData.OfficeName, officeData.EmailID, officeData.ContactNumber,
			officeData.WorkingHoursFrom, officeData.WorkingHoursTo, officeData.DivisionID, officeData.RegionID,
			officeData.CircleID, officeData.ReportingOfficeID, officeData.Latitude, officeData.Longitude,
			officeData.Status, officeData.CSIFacilityID, officeData.OpenToPublicDate, officeData.ClosedDate,
			officeData.ReasonForDisable, officeData.ReasonToEnable, officeData.CreatedBy, officeData.UpdatedBy,
			officeData.UpdatedDate, officeData.ValidatedFlag, officeID)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data in the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Office updated successfully"})
	}
}

func updateOfficeAttributeHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		attributeIDStr := c.Param("AttributeID")
		attributeID, err := strconv.Atoi(attributeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid AttributeID"})
			return
		}

		var officeAttributeData OfficeAttributeData
		if err := c.ShouldBindJSON(&officeAttributeData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON request"})
			return
		}

		_, err = db.Exec(`
            UPDATE OfficeAttributeMaster
            SET OfficeTypeID = $1, OpenedDate = $2, ClosedDate = $3, QRTerminalID = $4, OfficeAddressLine1 = $5, OfficeAddressLine2 = $6,
                OfficeAddressLine3 = $7, Landmark = $8, CityID = $9, DistrictID = $10, TalukID = $11, VillageID = $12, StateID = $13, Pincode = $14,
                PAOCode = $15, SolId = $16, PLIId = $17, GSTNForHO = $18, WEGCode = $19, DDOCode = $20, DeliveryOfficeFlag = $21,
                CSIRolledOutFlag = $22, SingleHandedOfficeFlag = $23, CreatedBy = $24, CreatedDate = $25, UpdatedBy = $26, UpdatedDate = $27
            WHERE AttributeID = $28`,
			officeAttributeData.OfficeTypeID, officeAttributeData.OpenedDate, officeAttributeData.ClosedDate,
			officeAttributeData.QRTerminalID, officeAttributeData.OfficeAddressLine1, officeAttributeData.OfficeAddressLine2,
			officeAttributeData.OfficeAddressLine3, officeAttributeData.Landmark, officeAttributeData.CityID,
			officeAttributeData.DistrictID, officeAttributeData.TalukID, officeAttributeData.VillageID,
			officeAttributeData.StateID, officeAttributeData.Pincode, officeAttributeData.PAOCode, officeAttributeData.SolId,
			officeAttributeData.PLIId, officeAttributeData.GSTNForHO, officeAttributeData.WEGCode, officeAttributeData.DDOCode,
			officeAttributeData.DeliveryOfficeFlag, officeAttributeData.CSIRolledOutFlag, officeAttributeData.SingleHandedOfficeFlag,
			officeAttributeData.CreatedBy, officeAttributeData.CreatedDate, officeAttributeData.UpdatedBy,
			officeAttributeData.UpdatedDate, attributeID)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data in the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Office attribute updated successfully"})
	}
}
