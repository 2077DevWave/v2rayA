#!/bin/sh

systemctl daemon-reload
systemctl enable v2raya
systemctl start v2raya
echo -e "\033[36m************************************\033[0m"
echo -e "\033[36m*        Congratulations!          *\033[0m"
echo -e "\033[36m* GUI demo: https://v2raya.mzz.pub *\033[0m"
echo -e "\033[36m************************************\033[0m"
