select *
from mgo.products;

SELECT *
FROM mgo.products
where name = 'Servilletas Elite 420 pz';

select * from mgo.categories_subcategories
where categories_subcategories.sub_categories_id in ('64e7c19bdf724c5adae20abf','64e7c19bbd79ca5ada7b92d1','5b56bac4f6c3f0d9b33f1784','64e7c19b8a15db5adae09c0c')
  and categories_subcategories.categories_id in ('64e7c19bdf724c5adae20abf','64e7c19bbd79ca5ada7b92d1','5b56bac4f6c3f0d9b33f1784','64e7c19b8a15db5adae09c0c');


SELECT DISTINCT ON (p.sku) p.sku,
                           p.barcode,
                           p.name,
                           p.short_description,
                           p.variant,
                           (select min(price_list) as min_price from mgo.products_inventory where sku = p.sku and zone is not null),
                           (select max(price_list) as max_price from mgo.products_inventory where sku = p.sku and zone is not null),
                           b.name  as brand_name,
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

select count(*)
from mgo.products;

select *
from mgo.categories
where _id in ('64e7c19bdf724c5adae20abf','64e7c19bbd79ca5ada7b92d1','5b56bac4f6c3f0d9b33f1784','64e7c19b8a15db5adae09c0c');

select * from mgo.categories_subcategories
where categories_subcategories.sub_categories_id in ('64e7c19bdf724c5adae20abf','64e7c19bbd79ca5ada7b92d1','64e7c19b8a15db5adae09c0c')
    and categories_subcategories.categories_id = '5b56bac4f6c3f0d9b33f1784';

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

select distinct c._id as id, c.name as name, s._id as subcategory_id, s.name as subcategory_name, c.image as image
from mgo.categories c
         join mgo.categories_subcategories cs on c._id = cs.categories_Id
         join mgo.subcategories s on cs.sub_categories_id = s._id
union all
select c._id as id, c.name as name, null as subcategory_id, null as subcategory_name, c.image as image
from mgo.categories c
order by id, subcategory_id asc;


select * from mgo.products where name = 'Papel Aluminio Reynolds Wrap 5m x 30cm 1 pz';

select *
from mgo.products_inventory;