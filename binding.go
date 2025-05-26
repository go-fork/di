package di

// BindingFunc là một hàm trả về một instance của dependency.
// BindingFunc được sử dụng để đăng ký các dependency trong container.
type BindingFunc func(c *Container) interface{}
