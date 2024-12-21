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
	r.rating
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


WITH WeightedAverages AS (
    SELECT 
        r.ticket_id,
        SUM(r.rating * c.weight) / SUM(c.weight) AS weighted_average
    FROM 
        ratings r
    JOIN 
        rating_categories c ON r.rating_category_id = c.id
    GROUP BY 
        r.ticket_id
)
SELECT 
    ticket_id,
    ROUND(weighted_average / 5 * 100, 2) AS score 
FROM 
    WeightedAverages;


with normalized_ticket_rating as (
		select t.id,
			r.rating_category_id,
			((r.rating / 5.0) * 100) as normalized_rating
		from tickets t
			left join ratings r on (t.id = r.ticket_id)
			join rating_categories rc on (r.rating_category_id = rc.id)
			-- and r.created_at between '2020-01-01T00:00:00' AND '2020-01-31T23:59:00'
	)
select id,
	rating_category_id,
	avg(normalized_rating)
from normalized_ticket_rating
where id=88
group by id,
	rating_category_id



WITH FilteredRatings AS (
    SELECT
        t.id as ticket_id,
        r.rating_category_id,
        r.rating,
        c.weight
    FROM
        ratings r
    JOIN
        rating_categories c ON r.rating_category_id = c.id
    JOIN
        tickets t ON r.ticket_id = t.id
    WHERE
        r.created_at BETWEEN '2019-01-01T00:00:00' AND '2020-05-02T00:00:00' -- Replace with your desired date range
),
WeightedAverages AS (
    SELECT
        ticket_id,
        rating_category_id,
        SUM(rating * weight) / SUM(weight) AS weighted_average
    FROM
        FilteredRatings
    GROUP BY
        ticket_id,
        rating_category_id
)
SELECT
    ticket_id,
    rating_category_id,
    ROUND(weighted_average / 5 * 100, 2) AS category_score
FROM
    WeightedAverages WHERE ticket_id=60695;


select * from ratings where ticket_id = 60695;

select ticket_id,rating_category_id from ratings where rating_category_id=1 and created_at BETWEEN '2019-03-03T00:00:00' AND '2019-03-10T00:00:00'



WITH FilteredRatings AS (
    SELECT
        t.id as ticket_id,
        r.rating_category_id,
        r.rating,
        c.weight,
        r.created_at as review_date
    FROM
        ratings r
    JOIN
        rating_categories c ON r.rating_category_id = c.id
    JOIN
        tickets t ON r.ticket_id = t.id
    WHERE
        r.created_at BETWEEN  '2019-01-01T00:00:00' AND '2020-05-02T00:00:00' 
),
DailyAverages AS (
    SELECT
        rating_category_id,
        date(review_date) AS review_date,
        AVG(rating) AS daily_average_rating,
        COUNT(*) AS daily_rating_count
    FROM
        FilteredRatings
    GROUP BY
        rating_category_id,
        date(review_date)
)
,
WeeklyAverages AS (
    SELECT
        rating_category_id,
        date(review_date, 'weekday 1') AS week_start,
        date(review_date, 'weekday 1', '+6 days') AS week_end,
        AVG(daily_average_rating) AS weekly_average_rating,
        SUM(daily_rating_count) AS weekly_rating_count
    FROM
        DailyAverages
    GROUP BY
        rating_category_id,
        date(review_date, 'weekday 0')
)
,
AggregatedScores AS (
    SELECT
        DailyAverages.rating_category_id,
        CASE 
            WHEN (julianday('2020-05-02T00:00:00' ) - julianday('2019-01-01T00:00:00')) <= 31 THEN review_date 
            ELSE week_start || ' - ' || week_end 
        END AS aggregation_period,
        CASE 
            WHEN (julianday('2020-05-02T00:00:00' ) - julianday('2019-01-01T00:00:00')) <= 31 THEN daily_average_rating 
            ELSE weekly_average_rating 
        END AS average_rating,
        CASE 
            WHEN (julianday('2020-05-02T00:00:00') - julianday('2019-01-01T00:00:00')) <= 31 THEN daily_rating_count 
            ELSE weekly_rating_count 
        END AS rating_count
    FROM
        DailyAverages
    LEFT JOIN
        WeeklyAverages ON DailyAverages.rating_category_id = WeeklyAverages.rating_category_id
            AND DailyAverages.review_date = WeeklyAverages.week_start
    UNION ALL
    SELECT
        WeeklyAverages.rating_category_id,
        CASE 
            WHEN (julianday('2020-05-02T00:00:00' ) - julianday('2019-01-01T00:00:00')) <= 31 THEN review_date 
            ELSE week_start || ' - ' || week_end 
        END AS aggregation_period,
        CASE 
            WHEN (julianday('2020-05-02T00:00:00' ) - julianday('2019-01-01T00:00:00')) <= 31 THEN daily_average_rating 
            ELSE weekly_average_rating 
        END AS average_rating,
        CASE 
            WHEN (julianday('2020-05-02T00:00:00') - julianday('2019-01-01T00:00:00')) <= 31 THEN daily_rating_count 
            ELSE weekly_rating_count 
        END AS rating_count
    FROM
        WeeklyAverages
    LEFT JOIN
        DailyAverages ON WeeklyAverages.rating_category_id = DailyAverages.rating_category_id
            AND WeeklyAverages.week_start = DailyAverages.review_date
    WHERE
        DailyAverages.review_date IS NULL
)
SELECT
    rating_category_id,
    aggregation_period,
    ROUND(average_rating / 5 * 100, 2) AS category_score,
    rating_count
FROM
    AggregatedScores
WHERE aggregation_period IS NOT NULL
ORDER BY
    rating_category_id,
    aggregation_period;



WITH FilteredRatings AS (
    SELECT
        t.id as ticket_id,
        r.rating_category_id,
        c.name as rating_category_name,
        r.rating,
        c.weight,
        r.created_at as review_date
    FROM
        ratings r
    JOIN
        rating_categories c ON r.rating_category_id = c.id
    JOIN
        tickets t ON r.ticket_id = t.id
    WHERE
        r.created_at BETWEEN  '2019-03-01T00:00:00' AND '2019-03-30T00:00:00' 
),
DailyAverages AS (
    SELECT
        rating_category_id,
        rating_category_name,
        date(review_date) AS review_date,
        AVG(rating) AS daily_average_rating,
        COUNT(*) AS daily_rating_count
    FROM
        FilteredRatings
    GROUP BY
        rating_category_id,
        rating_category_name,
        date(review_date)
)
,
WeeklyAverages AS (
    SELECT
        rating_category_id,
        rating_category_name,
        date(review_date, 'weekday 1') AS week_start,
        date(review_date, 'weekday 1', '+6 days') AS week_end,
        AVG(daily_average_rating) AS weekly_average_rating,
        SUM(daily_rating_count) AS weekly_rating_count
    FROM
        DailyAverages
    GROUP BY
        rating_category_id,
        rating_category_name,
        date(review_date, 'weekday 0')
)
,
AggregatedScores AS (
    SELECT
        DailyAverages.rating_category_id,
        CASE 
            WHEN (julianday('2019-03-30T00:00:00') - julianday('2019-03-01T00:00:00')) <= 31 THEN review_date 
            ELSE week_start || ' - ' || week_end 
        END AS aggregation_period,
        CASE 
            WHEN (julianday('2019-03-30T00:00:00' ) - julianday('2019-03-01T00:00:00')) <= 31 THEN daily_average_rating 
            ELSE weekly_average_rating 
        END AS average_rating,
        CASE 
            WHEN (julianday('2019-03-30T00:00:00') - julianday('2019-03-01T00:00:00')) <= 31 THEN daily_rating_count 
            ELSE weekly_rating_count 
        END AS rating_count
    FROM
        DailyAverages
    LEFT JOIN
        WeeklyAverages ON DailyAverages.rating_category_id = WeeklyAverages.rating_category_id
            AND DailyAverages.review_date = WeeklyAverages.week_start
    UNION ALL
    SELECT
        WeeklyAverages.rating_category_id,
        CASE 
            WHEN (julianday('2019-03-30T00:00:00' ) - julianday('2019-03-01T00:00:00')) <= 31 THEN review_date 
            ELSE week_start || ' - ' || week_end 
        END AS aggregation_period,
        CASE 
            WHEN (julianday('2019-03-30T00:00:00' ) - julianday('2019-03-01T00:00:00')) <= 31 THEN daily_average_rating 
            ELSE weekly_average_rating 
        END AS average_rating,
        CASE 
            WHEN (julianday('2019-03-30T00:00:00') - julianday('2019-03-01T00:00:00')) <= 31 THEN daily_rating_count 
            ELSE weekly_rating_count 
        END AS rating_count
    FROM
        WeeklyAverages
    LEFT JOIN
        DailyAverages ON WeeklyAverages.rating_category_id = DailyAverages.rating_category_id
            AND WeeklyAverages.week_start = DailyAverages.review_date
    WHERE
        DailyAverages.review_date IS NULL
)
SELECT
    rating_category_id,
    aggregation_period,
    ROUND(average_rating / 5 * 100, 2) AS category_score,
    rating_count
FROM
    AggregatedScores
WHERE aggregation_period IS NOT NULL
ORDER BY
    rating_category_id,
    aggregation_period;



WITH FilteredRatings AS (
    SELECT
        r.rating,
        c.weight
    FROM
        ratings r
    JOIN
        rating_categories c ON r.rating_category_id = c.id
    WHERE
        r.created_at BETWEEN '2019-03-01T00:00:00' AND '2019-03-30T00:00:00' 
),
WeightedAverage AS (
    SELECT 
        SUM(rating * weight) / SUM(weight) AS overall_average_rating
    FROM 
        FilteredRatings r
)
SELECT 
    'Overall Score' AS period,
    ROUND(AVG(overall_average_rating) / 5 * 100, 2) AS overall_score 
FROM 
    WeightedAverage;