-- with normalized_rating_categories as (
-- 	SELECT
--         id,
-- 		weight,
--         SUM(weight) OVER () AS normalized_weight
--     FROM
--         rating_categories
-- )
select t.id,
	r.rating_category_id,
	r.created_at,
	rc.weight,
	r.rating,
	(r.rating / 5.0) as normalized_rating,
	(
		((r.rating / 5.0) * rc.weight) / SUM(rc.weight)
	) * 100
from tickets t
	left join ratings r on (t.id = r.ticket_id)
	inner join rating_categories rc on (r.rating_category_id = rc.id)
where t.id = 60695 


with normalized_ticket_rating as (
		select t.id,
			r.rating_category_id,
			((r.rating / 5.0) * 100) as normalized_rating
		from tickets t
			left join ratings r on (t.id = r.ticket_id)
			join rating_categories rc on (r.rating_category_id = rc.id)
			and r.created_at between '2020-01-01T00:00:00' AND '2020-01-31T23:59:00'
	)
select id,
	rating_category_id,
	avg(normalized_rating)
from normalized_ticket_rating
group by id,
	rating_category_id

-- rating / 5 = rating individual

-- (rating individual * weight category) / sum(weight)