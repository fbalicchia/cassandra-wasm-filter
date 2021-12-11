#/bin/bash

# create tables

docker exec cass-server sh -c "cqlsh -e \"CREATE KEYSPACE deathstar WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : 2 };\""
docker exec cass-server sh -c "cqlsh -e \"CREATE TABLE deathstar.scrum_notes (creation timeuuid, empire_member_id uuid, content varchar, PRIMARY KEY (empire_member_id));\"" 

UUID1=`uuidgen`
UPDATE1="Trying to figure out if we should paint it medium grey, light grey, or medium-light grey.  Not blocked."
QUERY1="INSERT INTO deathstar.scrum_notes (empire_member_id, creation, content) values ($UUID1, now(), '${UPDATE1}');"
UUID2=`uuidgen`
UPDATE2="I think the exhaust port could be vulnerable to a direct hit.  Hope no one finds out about it.  Not blocked."
QUERY2="INSERT INTO deathstar.scrum_notes (empire_member_id, creation, content) values ($UUID2, now(), '${UPDATE2}');"
UUID3=`uuidgen`
UPDATE3="Designed protective shield for deathstar.  Could be based on nearby moon.  Feature punted to v2.  Not blocked."
QUERY3="INSERT INTO deathstar.scrum_notes (empire_member_id, creation, content) values ($UUID3, now(), '${UPDATE3}');"

docker exec cass-server sh -c "cqlsh -e \"$QUERY1\""
docker exec cass-server sh -c "cqlsh -e \"$QUERY2\""
docker exec cass-server sh -c "cqlsh -e \"$QUERY3\""

