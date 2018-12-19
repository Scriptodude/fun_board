# Little hack to use relative paths in the go project
# I know it is not idiom go-lang, but it is much easier to follow imo.
ln -s $(pwd)/server/ $GOPATH/src/server
mkdir -p bin/static
cd server
go build
mv server ../bin/
cd ..
cp -r client/src/* bin/static

rm -rf $GOPATH/src/server