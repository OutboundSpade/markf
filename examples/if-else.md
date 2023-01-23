#!(var age = 99)
#!(if #!(var age) >= 18 {
	You are an adult!
} else {
	#!(if #!(var age) < 13 {
		You are a child!
	} else {
		You are a teenager!
	})
})


#!(if #!(var age) < 18 {
	You are not an adult!
})