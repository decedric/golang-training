#!/bin/bash

print_output_and_store_id () {
  if [ "$2" = "$(jq -sr '.[1].code' <<< "$1")" ];
  then
    echo "Success"
    id=$(jq -sr '.[0].id' <<<"$1")
    jq <<< "$1"
  else
    echo "Unexpected Return Type"
    echo "$1"
  fi
}

print_output () {
  local item="$(jq -sr '.[1].code' <<< "$1")"
  if [[ "$2" =~ "$item" ]];
  then
    echo "Success"
    jq <<< "$1"
  else
    echo "Unexpected Return Type:"
    echo "$1"
  fi
}
select action in start_fibonacci_100 poll_last get_result_last start_fibonacci_choose_number poll_fibonacci_choose_id get_result_fibonacci_choose_id
do
    case $action in
        start_fibonacci_100)
            content=$(curl -X POST -w '{"code": "%{http_code}"}' localhost:8080/fibonacci/100)
            print_output_and_store_id "$content" "201"
            ;;
        poll_last)
            content=$(curl -X GET -w '{"code": "%{http_code}"}' localhost:8080/fibonacci/polling/$id)
            print_output "$content" "200 302"
            ;;
        get_result_last)
            content=$(curl -X GET -w '{"code": "%{http_code}"}' localhost:8080/fibonacci/$id)
            print_output "$content" "200"
            ;;
        start_fibonacci_choose_number)
            echo "enter fibonacci position to compute: "
            read value
            content=$(curl -X POST -w '{"code": "%{http_code}"}' localhost:8080/fibonacci/$value)
            print_output "$content" "201"
            ;;
        poll_fibonacci_choose_id)
            echo "enter polling id: "
            read value
            content=$(curl -X GET -w '{"code": "%{http_code}"}' localhost:8080/fibonacci/polling/$value)
            print_output "$content" "200 302"
            ;;
        get_result_fibonacci_choose_id)
            echo "enter id: "
            read value
            content=$(curl -X GET -w '{"code": "%{http_code}"}' localhost:8080/fibonacci/$value)
            print_output "$content" "200"
            ;;
        *)
           echo "Ooops";;
    esac
done

