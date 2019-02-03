package main

const (
schema = "CREATE TABLE `shelters` (" +
  "`inspire_id` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL," +
  "`begin_life` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`end_life` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`type_of_occupation` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`filter_type` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`serviceLev` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`service_type` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`pointOfCon` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`relatedPar` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`slots` int(11) DEFAULT NULL," +
  "`estate_id` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`shelter_id` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL," +
  "`address` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`city` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`municipality` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL," +
  "`position_lat` decimal(15,13) DEFAULT NULL," +
  "`position_long` decimal(15,13) DEFAULT NULL," +
  "PRIMARY KEY (`shelter_id`), " +
  "KEY `position_long` (`position_long`)," +
  "KEY `position_lat` (`position_lat`)" +
") ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;"
)
