package app

// func TestIsErrReplace(t *testing.T) {
// 	utests := []struct {
// 		scenario     string
// 		err          error
// 		isErrReplace bool
// 	}{
// 		{
// 			scenario:     "error is a replace error",
// 			err:          errors.New("test").Tag("replace", true),
// 			isErrReplace: true,
// 		},
// 		{
// 			scenario:     "error is not a replace error",
// 			err:          errors.New("test").Tag("test", true),
// 			isErrReplace: false,
// 		},
// 		{
// 			scenario:     "standard error is not a replace error",
// 			err:          fmt.Errorf("test"),
// 			isErrReplace: false,
// 		},
// 		{
// 			scenario:     "nil error is not a replace error",
// 			err:          nil,
// 			isErrReplace: false,
// 		},
// 	}

// 	for _, u := range utests {
// 		t.Run(u.scenario, func(t *testing.T) {
// 			res := isErrReplace(u.err)
// 			require.Equal(t, u.isErrReplace, res)
// 		})
// 	}
// }

// type mountTest struct {
// 	scenario string
// 	node     UI
// }

// func testMountDismount(t *testing.T, utests []mountTest) {
// 	for _, u := range utests {
// 		t.Run(u.scenario, func(t *testing.T) {
// 			testSkipNonWasm(t)

// 			n := u.node
// 			err := mount(n)
// 			require.NoError(t, err)
// 			testMounted(t, n)

// 			dismount(u.node)
// 			testDismounted(t, n)
// 		})
// 	}
// }

// func testMounted(t *testing.T, n UI) {
// 	require.NotNil(t, n.JSValue())
// 	require.True(t, n.Mounted())

// 	// switch n.Kind() {
// 	// case HTML, Component:
// 	// 	require.NotNil(t, n.self())
// 	// }

// 	for _, c := range n.children() {
// 		require.Equal(t, n, c.parent())
// 		testMounted(t, c)
// 	}
// }

// func testDismounted(t *testing.T, n UI) {
// 	require.Nil(t, n.JSValue())
// 	require.False(t, n.Mounted())

// 	// switch n.Kind() {
// 	// case HTML, Component:
// 	// 	require.Nil(t, n.self())
// 	// }

// 	for _, c := range n.children() {
// 		testDismounted(t, c)
// 	}
// }

// type updateTest struct {
// 	scenario   string
// 	a          UI
// 	b          UI
// 	matches    []TestUIDescriptor
// 	replaceErr bool
// }

// func testUpdate(t *testing.T, utests []updateTest) {
// 	for _, u := range utests {
// 		t.Run(u.scenario, func(t *testing.T) {
// 			testSkipNonWasm(t)

// 			err := mount(u.a)
// 			require.NoError(t, err)
// 			defer dismount(u.a)

// 			err = update(u.a, u.b)
// 			if u.replaceErr {
// 				require.Error(t, err)
// 				require.True(t, isErrReplace(err))
// 				return
// 			}

// 			require.NoError(t, err)

// 			for _, d := range u.matches {
// 				require.NoError(t, TestMatch(u.a, d))
// 			}
// 		})
// 	}
// }
