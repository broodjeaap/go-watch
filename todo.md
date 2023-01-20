# Todo
- comments
- safe escape {{ }} for pages
- 'jitter' for cronjobs after first start ?
- SLOW SQL
    - SELECT * FROM `filter_connections` WHERE watch_id = 1
    - SELECT * FROM `watches` WHERE `watches`.`id` = 7 ORDER BY `watches`.`id` LIMIT 1
    - SELECT * FROM `filters` WHERE watch_id = 1
    - INSERT INTO `filter_outputs` (`watch_id`,`name`,`value`,`time`) VALUES (7,"Caption","D","2023-01-│>
20 14:20:25.454") RETURNING `id`