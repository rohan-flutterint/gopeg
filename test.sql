with a as select data.id as id, distinct_count(data.version) as cnt
          from (select parse_json(data) as data
                from file("data.jsonl")
                where line % 2 == 0)
          where data.sum > 0
          group by data.id;

with b as select parse_json(data) as data from file("b.jsonl");

select *
from a
         join b on a.x = b.y;
