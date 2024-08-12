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
