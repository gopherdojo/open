module github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima

go 1.16

replace (
	github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/download => ./download
	github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/getheader => ./getheader
	github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/listen => ./listen
	github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/request => ./request
)

require (
	github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/download v0.0.0-00010101000000-000000000000
	github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/getheader v0.0.0-00010101000000-000000000000
	github.com/MizushimaToshihiko/gopherdojo-studyroom/kadai3-2/Mizushima/request v0.0.0-00010101000000-000000000000
	github.com/jessevdk/go-flags v1.5.0
	github.com/pkg/errors v0.9.1
)
