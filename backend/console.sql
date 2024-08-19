--GETS JOB DATA BASED ON JOB NUMBER
SELECT
    cm.id AS cut_material_id,
    cm.job,
    cm.material_code AS cut_material_material_code,
    cm.quantity AS cut_material_quantity,
    cm.stock_length,
    cm.length AS cut_material_length,
    p.id AS part_id,
    p.part_number,
    p.material_code AS part_material_code,
    p.length AS part_length,
    cmp.part_qty
FROM
    cut_materials cm
        JOIN
    cut_material_parts cmp ON cm.id = cmp.cut_material_id
        JOIN
    parts p ON cmp.part_id = p.id
WHERE
    cm.job = 'TEST';

--Get Total footage and qty of each length of each code
SELECT
    material_code,
    stock_length,
    SUM(quantity) AS total_quantity,
    stock_length * SUM(quantity) AS total_length
FROM
    cut_materials
WHERE
    job = 'TEST'
GROUP BY
    material_code, stock_length;

--SELECTS *
SELECT * FROM cut_material_parts;
SELECT * FROM parts;
SELECT * FROM jobs;


SELECT * FROM jobs;