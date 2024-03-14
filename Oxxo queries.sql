select *
from mgo.products;

SELECT sku,
       barcode,
       name,
       short_description,
       variant,
       brand,
       categories,
       image,
       tags,
       taxes,
       is_enabled
FROM mgo.products
order by sku
limit 10;

SELECT DISTINCT ON (p.sku) p.sku,
                           p.barcode,
                           p.name,
                           p.short_description,
                           p.variant,
                           b.name as brand_name,
                           p.categories,
                           c.name  as category_name,
                           sc.name as subcategory_name,
                           p.image,
                           p.tags,
                           p.taxes,
                           p.age_restriction,
                           p.restricted,
                           p.is_enabled
FROM mgo.products p
         LEFT JOIN mgo.brands b ON p.brand = b._id
         LEFT JOIN LATERAL unnest(p.categories) AS cat(category_id) on true
         JOIN mgo.categories c ON c._id = cat.category_id
         LEFT JOIN LATERAL unnest(p.categories) AS subcat(category_id) on true
         LEFT JOIN mgo.subcategories sc ON sc._id = subcat.category_id
order by p.sku, category_name, subcategory_name;

select count(*) from mgo.products;

select *
from mgo.categories
where _id = '5b56bac4f6c3f0d9b33f1783';

select *
from mgo.categories
where _id in (select sub_categories_id from mgo.categories_subcategories);


select *
from mgo.subcategories;

select *
from mgo.categories;

SELECT *
FROM mgo.products
order by sku;

SELECT sku, categories
FROM mgo.products
order by sku;


SELECT unnest(string_to_array('It s an example sentence.', ' ')) AS parts;

select distinct c._id as id, c.name as name, c.image as image, s._id as subcategory_id, s.name as subcategory_name
from mgo.categories c
         join mgo.categories_subcategories cs on c._id = cs.categories_Id
         join mgo.subcategories s on cs.sub_categories_id = s._id
order by c._id, s._id asc;


select *
from mgo.products_inventory;


-- Brand owner for Brands

select count(*)
from mgo.categories
where _id not in (select categories_id from mgo.categories_subcategories);
select count(*)
from mgo.subcategories
where _id not in (select sub_categories_id from mgo.categories_subcategories);


select *
from mgo.subcategories
where _id in
      ('64dd197ab5d2485b1d2bfdc7', '5b56bac4f6c3f0d9b33f1783', '64e7c19b61c0be5adaedc704', '64e7c19be8e5985ada2105e1');


select *
from mgo.categories_subcategories
where sub_categories_id in
      ('64e7c19b8a15db5adae09c0c', '64e7c19bd870e05adacfbc09', '64e7c19b703ece5ada000c43', '64e7c19be8e5985ada2105e1');