set -oue

echo '------------------------'
echo $USER
echo '------------------------'

go mod init tempfile

# Add go symlink to curent dir
ln -s $(which go) go
chmod a+rx go
chmod a+rx $(which go)

exit 
