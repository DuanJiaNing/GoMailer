function registerApp()
{
    echo "Sending request... $1"
    info=`curl -H "Content-Type:application/json" -X POST --data @$2 $1`
    echo $info
}

registerApp $1 $2