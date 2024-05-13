-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW _solvedquizzes AS
(
    select q.id,
           q.user_answer,
           q.type,
           q.created_at,
           q.starred,
           qm.media,
           ql.lesson,
           q.product_id,
           q.lesson_id,
           q.user_id,
           q.quiz_id,
           json_build_object(
               'first_name', u.first_name,
               'last_name', u.last_name,
               'avatar', u.avatar
           ) as author

    from public.quiz_solved as q

    -- media
    left join lateral (
        select
            json_agg(
                json_build_object(
                    'url', m.url,
                    'sources', m.sources,
                    'type', m.type
                )
            ) as media
        from public.related_media as rm
        join public.media as m on m.id = rm.media_id
        where rm.related_type = 'solved_quiz' and rm.related_id = q.id
    ) as qm on true

    -- lesson
    left join lateral (
        select
            json_build_object(
                'slug', l.slug,
                'name', l.name,
                'product', json_build_object(
                    'name', p.name,
                    'slug', p.slug
                )
            ) as lesson
        from public.lesson as l
        join public.product as p on p.id = l.product_id
        where l.id = q.lesson_id
    ) as ql on true

    -- author
    join public.user as u on u.id = q.user_id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW _solvedquizzes;
-- +goose StatementEnd
