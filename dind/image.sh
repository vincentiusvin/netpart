#! /bin/sh

while ! docker version > /dev/null; do
    echo "Waiting..."
    sleep 3
done

if [ -e "/images/$POSTGRES_IMAGE.tar" ]
then
    echo "Found image file, loading..."
    for file in /images/*.tar; do
        docker load < $file
    done
else
    echo "Did not find image file, pulling..."
    docker image pull $POSTGRES_IMAGE
    docker image save $POSTGRES_IMAGE > /images/$POSTGRES_IMAGE.tar
fi
