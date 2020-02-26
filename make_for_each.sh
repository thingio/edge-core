#!/bin/bash
cmd=$1
for dir in `ls -d */`;do
    echo " = = = = = = = = = = Current directory is \"${dir}\" = = = = = = = = = ="
    make WORK_DIR=${dir} ${cmd}
done