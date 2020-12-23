ARCH="$(grep ^otelcontribcol-all-sys Makefile|fmt -w 1|tail -n +2)"
ARCH=(${ARCH//|/ })
MATRIX="{\"include\":["
for i in "${!ARCH[@]}"; do
if ((i == 0)); then
    MATRIX+="{\"arch\":\"${ARCH[$i]}\"}"
else
    MATRIX+=",{\"arch\":\"${ARCH[$i]}\"}"
fi
done
MATRIX+="]}"
echo "::set-output name=matrix::$MATRIX"
