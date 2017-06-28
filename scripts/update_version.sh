# Get the version.
version=`cat VERSION`

# Write out the package.
cat << EOF > version.go
package littlefly

//go:generate bash ../scripts/update_version.sh

// Version ...
var Version = "$version"
EOF
