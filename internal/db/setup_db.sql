-- Create a database file (use sqlite3 command line to execute this script)

-- Create the parts table
CREATE TABLE IF NOT EXISTS parts (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     part_number TEXT NOT NULL UNIQUE,
                                     material_code TEXT NOT NULL,
                                     length REAL NOT NULL
);

-- Create the cut_materials table
CREATE TABLE IF NOT EXISTS cut_materials (
                                             id INTEGER PRIMARY KEY AUTOINCREMENT,
                                             job TEXT NOT NULL ,
                                             material_code TEXT NOT NULL,
                                             quantity INTEGER NOT NULL,
                                             stock_length REAL NOT NULL,
                                             length REAL NOT NULL
);

-- Create the cut_material_parts table
CREATE TABLE IF NOT EXISTS cut_material_parts (
                                                  cut_material_id INTEGER,
                                                  part_id INTEGER,
                                                  part_qty INTEGER,
                                                  PRIMARY KEY (cut_material_id, part_id),
                                                  FOREIGN KEY (cut_material_id) REFERENCES cut_materials(id),
                                                  FOREIGN KEY (part_id) REFERENCES parts(id)
);

CREATE INDEX IF NOT EXISTS idx_cut_materials_job ON cut_materials(job);



-- -- Create a database file (use sqlite3 command line to execute this script)
--
-- -- Create the parts table
-- CREATE TABLE IF NOT EXISTS parts (
--                                      id INTEGER PRIMARY KEY AUTOINCREMENT,
--                                      part_number TEXT NOT NULL UNIQUE,
--                                      material_code TEXT NOT NULL,
--                                      length REAL NOT NULL CHECK(length > 0)
-- );
--
-- -- Create the cut_materials table
-- CREATE TABLE IF NOT EXISTS cut_materials (
--                                              id INTEGER PRIMARY KEY AUTOINCREMENT,
--                                              job TEXT NOT NULL,
--                                              material_code TEXT NOT NULL,
--                                              quantity INTEGER NOT NULL CHECK(quantity > 0),
--                                              stock_length REAL NOT NULL CHECK
--                                                  (stock_length > 0),
--                                              length REAL NOT NULL CHECK(length > 0)
-- );
--
-- -- Create the cut_material_parts table
-- CREATE TABLE IF NOT EXISTS cut_material_parts (
--                                                   cut_material_id INTEGER,
--                                                   part_id INTEGER,
--                                                   part_qty INTEGER NOT NULL CHECK(part_qty > 0),
--                                                   PRIMARY KEY (cut_material_id, part_id),
--                                                   FOREIGN KEY (cut_material_id) REFERENCES cut_materials(id) ON DELETE CASCADE,
--                                                   FOREIGN KEY (part_id) REFERENCES parts(id) ON DELETE CASCADE
-- );
--
-- -- Create indexes for faster queries
-- CREATE INDEX IF NOT EXISTS idx_parts_material_code ON parts(material_code);
-- CREATE INDEX IF NOT EXISTS idx_cut_materials_material_code ON cut_materials(material_code);

