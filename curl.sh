select action in start_fibonacci_100 poll_last get_result_last start_fibonacci_choose_number poll_fibonacci_choose_id get_result_fibonacci_choose_id
do
    case $action in
        start_fibonacci_100)
            content=$(curl -X POST localhost:8080/fibonacci/100)
            id=$(jq -r '.id' <<<"$content")
            jq <<< "$content"
            ;;
        poll_last)
            content=$(curl -X GET localhost:8080/fibonacci/polling/$id)
            jq <<< "$content"
            ;;
        get_result_last)
            content=$(curl -X GET localhost:8080/fibonacci/$id)
            jq <<< "$content"
            ;;
        start_fibonacci_choose_number)
            echo "enter fibonacci position to compute: "
            read value
            content=$(curl -X POST localhost:8080/fibonacci/$value)
            id=$(jq -r '.id' <<<"$content")
            jq <<< "$content"
            ;;
        poll_fibonacci_choose_id)
            echo "enter polling id: "
            read value
            content=$(curl -X GET localhost:8080/fibonacci/polling/$value)
            jq <<< "$content"
            ;;
        get_result_fibonacci_choose_id)
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

