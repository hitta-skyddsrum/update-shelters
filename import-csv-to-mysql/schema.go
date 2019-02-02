package main

const (
schema = "CREATE TABLE `shelters` (" +
  "`InspireID` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL," +
  "`beginLifes` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`endLifeSpa` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`typeOfOccu` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`resourceFi` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`serviceLev` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`serviceTyp` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`pointOfCon` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`relatedPar` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`numberOfOc` int(11) DEFAULT NULL," +
  "`additional` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`name` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`serviceLBA` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`serviceLBC` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`serviceLBM` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`Latitude` decimal(15,13) DEFAULT NULL," +
  "`Longitude` decimal(15,13) DEFAULT NULL," +
  "PRIMARY KEY (`InspireID`)" +
") ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;"
)
