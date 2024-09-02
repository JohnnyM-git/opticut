-- Create the jobs table
CREATE TABLE IF NOT EXISTS jobs (
                                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                                    job_number TEXT NOT NULL COLLATE NOCASE,
                                    customer TEXT NOT NULL COLLATE NOCASE,
                                    star INTEGER NOT NULL DEFAULT 0,
                                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create the parts table
CREATE TABLE IF NOT EXISTS parts (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     part_number TEXT NOT NULL UNIQUE,
                                     material_code TEXT NOT NULL,
                                     length REAL NOT NULL,
                                     cutting_operation TEXT
);

-- Create the cut_materials table
CREATE TABLE IF NOT EXISTS cut_materials (
                                             id INTEGER PRIMARY KEY AUTOINCREMENT,
                                             job TEXT NOT NULL COLLATE NOCASE,
                                             job_id INTEGER NOT NULL,
                                             material_code TEXT NOT NULL,
                                             quantity INTEGER NOT NULL,
                                             stock_length REAL NOT NULL,
                                             length REAL NOT NULL,
                                             cutting_operation TEXT,
                                             FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- Create the cut_material_parts table
CREATE TABLE IF NOT EXISTS cut_material_parts (
                                                  cut_material_id INTEGER,
                                                  part_id INTEGER,
                                                  part_qty INTEGER,
                                                  total_part_qty INTEGER,
                                                  job_id INTEGER,
                                                  length REAL,
                                                  part_cut_length REAL,
                                                  material_code TEXT,
                                                  cutting_operation TEXT,
                                                  PRIMARY KEY (cut_material_id, part_id),
                                                  FOREIGN KEY (cut_material_id) REFERENCES cut_materials(id),
                                                  FOREIGN KEY (part_id) REFERENCES parts(id),
                                                  FOREIGN KEY (job_id) REFERENCES jobs(id)
);


CREATE INDEX IF NOT EXISTS idx_cut_materials_job ON cut_materials(job);

-- Trigger to delete old entries when inserting new ones with the same job_number
CREATE TRIGGER IF NOT EXISTS delete_old_job
    BEFORE INSERT ON jobs
    FOR EACH ROW
BEGIN
    DELETE FROM cut_material_parts WHERE job_id = (SELECT id FROM jobs WHERE job_number = NEW.job_number);
    DELETE FROM cut_materials WHERE job_id = (SELECT id FROM jobs WHERE job_number = NEW.job_number);
    DELETE FROM jobs WHERE job_number = NEW.job_number;
END;

CREATE TRIGGER IF NOT EXISTS update_existing_part
    BEFORE INSERT ON parts
    FOR EACH ROW
BEGIN
    -- Check if a record with the same part_number already exists
    UPDATE parts
    SET length = NEW.length,
        material_code = NEW.material_code
    WHERE part_number = NEW.part_number;

    -- If the part_number does not exist, the record will be inserted normally
    -- The insertion operation will not be blocked by the trigger.
END;
