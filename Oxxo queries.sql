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
                           (select max(max_per_order) as max_quantity from mgo.products_inventory where sku = p.sku and zone is not null),
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
         WHERE (c.name, coalesce(sc.name, '')) in (select c.name, s.name
         from mgo.categories c
         join mgo.categories_subcategories cs on c._id = cs.categories_Id
         join mgo.subcategories s on cs.sub_categories_id = s._id
         union all
         select c.name as name, '' as subcategory_name from mgo.categories c where c.name != 'Promociones')
order by p.sku, category_name, subcategory_name;



select distinct c.name as name, s.name as subcategory_name
from mgo.categories c
         join mgo.categories_subcategories cs on c._id = cs.categories_Id
         join mgo.subcategories s on cs.sub_categories_id = s._id
union all
select c.name as name, null as subcategory_name
from mgo.categories c;


select * from mgo.products where name = 'Papel Aluminio Reynolds Wrap 5m x 30cm 1 pz';

select * from mgo.products where name = 'Pistaches Wonderful Pistachios 100 g';

select * from mgo.products where sku in ('7501048810307', '7501000111459', '21136020373');


select sku, store, zone, cr, price_list from mgo.products_inventory pi
                                                 LEFT JOIN mgo.products_inventory_stores pis on pi.store = pis._id


