package base62

import "testing"

func TestBase62Encode(t *testing.T) {
	type args struct{
		num int
	}
	cases := []struct{
		name string
		args args
		want string
	}{
		{
			name: "case1",
			args: args{
				num: 123453,
			},
			want: "w7b",
		},
	}

	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			if res := Base62Encode(ca.args.num); res != ca.want {
				t.Fatalf("Excepted: %v, but got: %v", ca.want, res)
			}
		})
	}
}

func TestBase62Decode(t *testing.T) {
	type args struct{
		s string
	}
	cases := []struct{
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{
				s: "w7b",
			},
			want: 123453,
		},
	}

	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			if res, err := Base62Decode(ca.args.s); err != nil || res != ca.want {
				t.Fatalf("Excepted: %v, but got: %v", ca.want, res)
			}
		})
	}
}