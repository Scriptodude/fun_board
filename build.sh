# Little hack to use relative paths in the go project
# I know it is not idiom go-lang, but it is much easier to follow imo.
ln -s $(pwd)/server/ $GOPATH/src/server
mkdir -p bin
mkdir -p /tmp/fun_board
cd server
go build
mv server ../bin/
cd ..
cp -r client/src/* /tmp/fun_board

rm -rf $GOPATH/src/server