select action in run poll end 
do
    case $action in
        run)
            curl -X POST localhost:8080/fibonacci/100
            echo "\n"
            ;;
        poll)
            echo "enter polling id: "
            read value
            echo "\n"
            curl -X GET localhost:8080/fibonacci/polling/$value
            echo "\n"
            ;;
        end)
            echo "enter id: "
            read value
            echo "\n"
            curl -X GET localhost:8080/fibonacci/$value
            echo "\n"
            ;;
        *)
           echo "Ooops";;
    esac
done
#!/bin/bash

