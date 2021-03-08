# Downloads BusyBox and creates symlinks for each of the 
# commands that it supports so that they are available
# by ther Linux common names (e.g. ls, hostname, et cetera)

BUSYBOX_ROOT=./bb_root/bin

echo "Creating BusyBox directory: $BUSYBOX_ROOT"
mkdir -p $BUSYBOX_ROOT
cd $BUSYBOX_ROOT

echo "Downloading BusyBox..."
curl https://www.busybox.net/downloads/binaries/1.30.0-i686/busybox > busybox
chmod u+x busybox

echo "Creating symlinks..."
for i in $(busybox --list)
do
    if [ "$i" != "busybox" ]
    then
        ln -s busybox $i
    fi
done

echo "Done"