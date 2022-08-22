select action in run poll end run1 poll1 end1
do
    case $action in
        run)
            content=$(curl -X POST localhost:8080/fibonacci/100)
            id=$(jq -r '.id' <<<"$content")
            jq <<< "$content"
            ;;
        poll)
            content=$(curl -X GET localhost:8080/fibonacci/polling/$id)
            jq <<< "$content"
            ;;
        end)
            content=$(curl -X GET localhost:8080/fibonacci/$id)
            jq <<< "$content"
            ;;
        run1)
            echo "enter fibonacci position to compute: "
            read value
            content=$(curl -X POST localhost:8080/fibonacci/$value)
            id=$(jq -r '.id' <<<"$content")
            jq <<< "$content"
            ;;

        poll1)
            echo "enter polling id: "
            read value
            content=$(curl -X GET localhost:8080/fibonacci/polling/$value)
            jq <<< "$content"
            ;;
        end1)
            echo "enter id: "
            read value
            content=$(curl -X GET localhost:8080/fibonacci/$value)
            jq <<< "$content"
            ;;
        *)
           echo "Ooops";;
    esac
done
#!/bin/bash

