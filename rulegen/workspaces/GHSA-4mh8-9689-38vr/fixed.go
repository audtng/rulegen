package main

# Disallow writes to the well-known directory included in
# the user's PATH on several distributions
audit deny @{HOME}/bin/{,**} wl,
audit deny @{HOME}/bin wl,
`

const homeConnectedPlugAppArmorWithAllRead = `
