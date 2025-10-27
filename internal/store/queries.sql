-- name: GetVideos :many
SELECT id, title, course_id, storage_key FROM videosections;

-- name: GetCourseMaterials :many
SELECT id, name, course_id, storage_key FROM course_materials;
