-- Scores by ticket,
-- not filtering,
-- including all tickets in the database
WITH ratings_with_score AS(
    SELECT r.ticket_id,
        r.rating_category_id,
        rc.name,
        COALESCE(
            ROUND(
                (r.rating * rc.weight * 1.0) / (5 * rc.weight) * 100,
                2
            ),
            0
        ) AS score
    FROM ratings r
        INNER JOIN rating_categories rc ON r.rating_category_id = rc.id
    WHERE r.created_at BETWEEN '2019-01-01T00:00:00' AND '2019-12-01T00:00:00'
),
ticket_categories AS(
    SELECT t.id AS ticket_id,
        rt.id AS category_id
    FROM tickets t
        CROSS JOIN rating_categories rt
)
SELECT tc.ticket_id,
    tc.category_id,
    COALESCE(score, 0) AS score
FROM ticket_categories tc
    LEFT JOIN ratings_with_score rs ON (
        tc.ticket_id = rs.ticket_id
        AND tc.category_id = rs.rating_category_id
    );
-- Scores by ticket,
-- filtering ticket that actually have rating 
WITH ratings_with_score AS(
    SELECT r.ticket_id,
        r.rating_category_id AS category_id,
        rc.name AS category_name,
        COALESCE(
            ROUND(
                (r.rating * rc.weight * 1.0) / (5 * rc.weight) * 100,
                2
            ),
            0
        ) AS score
    FROM ratings r
        INNER JOIN rating_categories rc ON r.rating_category_id = rc.id
    WHERE r.created_at BETWEEN '2019-01-01T00:00:00' AND '2019-12-01T00:00:00'
)
SELECT ticket_id,
    category_id,
    category_name,
    score
FROM ratings_with_score;

-- Aggregated category scores over a period of time
WITH RECURSIVE dates(date) AS (
  VALUES('2019-07-17')
  UNION ALL
  SELECT date(date, '+1 day')
  FROM dates
  WHERE date < '2019-07-30'
),
ratings_with_score AS (
    SELECT 
        r.ticket_id,
        r.rating_category_id AS category_id,
        rc.name AS category_name,
        r.rating,
        r.created_at,
        COALESCE(
            ROUND(
                (r.rating * rc.weight * 1.0) / (5 * rc.weight) * 100,
                2
            ),
            0
        ) AS score
    FROM ratings r
        LEFT JOIN rating_categories rc ON r.rating_category_id = rc.id
    WHERE r.created_at BETWEEN '2019-07-17T00:00:00' AND '2019-07-30T00:00:00'
)
SELECT rs.category_id, 
       COUNT(rs.rating) AS ratings, 
       rs.created_at
FROM ratings_with_score rs
LEFT JOIN dates dr 
    ON rs.created_at >= dr.date AND rs.created_at < dr.date
GROUP BY rs.category_id, dr.date;



WITH RECURSIVE dates(date) AS (
  VALUES('2019-07-17')
  UNION ALL
  SELECT date(date, '+1 day')
  FROM dates
  WHERE date < '2019-07-18'
)
SELECT date FROM dates;