package renderer

import (
	"VulpesEditor/util"
	_ "embed"
	"fmt"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const FLOAT_SIZE = 4

//go:embed shaders/triangle.vert
var vertexShader string

//go:embed shaders/triangle.frag
var fragmentShader string

//go:embed shaders/texture.vert
var texVertexShader string

//go:embed shaders/texture.frag
var texFragmentShader string

type uniforms struct {
	color int32
}

type renderer struct {
	shaderHandle uint32
	vao, vbo     uint32
	uniforms     uniforms
}

type textureUniforms struct {
	zoom    int32
	move    int32
	aspect  int32
	outline int32
	texUnit int32
	preUnit int32
	size    int32
}

type textureRender struct {
	shaderHandle uint32
	textureVao   uint32
	outlineVao   uint32
	vbo          uint32
	uniforms     textureUniforms
}

type windowScreen struct {
	width, height int32
}

type FrameBuffer struct {
	fbo           uint32
	colorBuffer   uint32
	depth         uint32
	width, height int32
}

var w windowScreen
var r renderer
var rTex textureRender

func glError(handle uint32, statusType uint32, getIV func(uint32, uint32, *int32), getInfoLog func(uint32, int32, *int32, *uint8), failureMsg string) {
	var status int32
	getIV(handle, statusType, &status)
	if status == gl.FALSE {
		var logLength int32
		getIV(handle, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength))
		getInfoLog(handle, logLength, nil, gl.Str(infoLog))
		fmt.Println(failureMsg+"\n", infoLog)
	}
}

func CreateTexture(width, height int32, data []float32) (id uint32) {
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16, width, height, 0, gl.RGBA, gl.FLOAT, gl.Ptr(data))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return id
}

func WriteTexture(id uint32, width, height int32, data []float32) {
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16, width, height, 0, gl.RGBA, gl.FLOAT, gl.Ptr(data))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func Init() {
	glShaderSource := func(handle uint32, source string) {
		csource, free := gl.Strs(source + "\x00")
		defer free()

		gl.ShaderSource(handle, 1, csource, nil)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

	w.width = 1200
	w.height = 900

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	{
		r.shaderHandle = gl.CreateProgram()
		vertHandle := gl.CreateShader(gl.VERTEX_SHADER)
		fragHandle := gl.CreateShader(gl.FRAGMENT_SHADER)
		glShaderSource(vertHandle, vertexShader)
		glShaderSource(fragHandle, fragmentShader)
		gl.CompileShader(vertHandle)
		glError(vertHandle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "Vertex shader error")
		gl.CompileShader(fragHandle)
		glError(vertHandle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "Fragment shader error")
		gl.AttachShader(r.shaderHandle, vertHandle)
		gl.AttachShader(r.shaderHandle, fragHandle)
		gl.LinkProgram(r.shaderHandle)
		glError(r.shaderHandle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog, "Linking program error")
		gl.DeleteShader(vertHandle)
		gl.DeleteShader(fragHandle)

		r.uniforms.color = gl.GetUniformLocation(r.shaderHandle, util.Str("color"))

		vertices := []float32{
			-0.5, -0.5, 0.0,
			0.5, -0.5, 0.0,
			0.0, 0.5, 0.0,
		}

		gl.GenVertexArrays(1, &r.vao)
		gl.GenBuffers(1, &r.vbo)
		gl.BindVertexArray(r.vao)

		gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, int(FLOAT_SIZE)*len(vertices), gl.Ptr(&vertices[0]), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*FLOAT_SIZE, nil)
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	{
		rTex.shaderHandle = gl.CreateProgram()
		vertHandle := gl.CreateShader(gl.VERTEX_SHADER)
		fragHandle := gl.CreateShader(gl.FRAGMENT_SHADER)
		glShaderSource(vertHandle, texVertexShader)
		glShaderSource(fragHandle, texFragmentShader)
		gl.CompileShader(vertHandle)
		glError(vertHandle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "Vertex shader error")
		gl.CompileShader(fragHandle)
		glError(vertHandle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "Fragment shader error")
		gl.AttachShader(rTex.shaderHandle, vertHandle)
		gl.AttachShader(rTex.shaderHandle, fragHandle)
		gl.LinkProgram(rTex.shaderHandle)
		glError(rTex.shaderHandle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog, "Linking program error")
		gl.DeleteShader(vertHandle)
		gl.DeleteShader(fragHandle)

		rTex.uniforms.move = gl.GetUniformLocation(rTex.shaderHandle, util.Str("move"))
		rTex.uniforms.zoom = gl.GetUniformLocation(rTex.shaderHandle, util.Str("zoom"))
		rTex.uniforms.aspect = gl.GetUniformLocation(rTex.shaderHandle, util.Str("aspect"))
		rTex.uniforms.outline = gl.GetUniformLocation(rTex.shaderHandle, util.Str("outline"))
		rTex.uniforms.texUnit = gl.GetUniformLocation(rTex.shaderHandle, util.Str("art"))
		rTex.uniforms.preUnit = gl.GetUniformLocation(rTex.shaderHandle, util.Str("preview"))
		rTex.uniforms.size = gl.GetUniformLocation(rTex.shaderHandle, util.Str("size"))
		gl.UseProgram(rTex.shaderHandle)
		gl.Uniform1i(rTex.uniforms.texUnit, 0)
		gl.Uniform1i(rTex.uniforms.preUnit, 1)

		textureVertices := []float32{
			1.0, -1.0, 1.0, 1.0,
			-1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, 0.0, 1.0,
			-1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0,
			1.0, 1.0, 1.0, 0.0,
		}

		outlineVertices := []float32{
			1.0, -1.0, 1.0, 1.0,
			-1.0, -1.0, 0.0, 1.0,
			-1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 0.0,
		}

		gl.GenVertexArrays(1, &rTex.textureVao)
		gl.GenBuffers(1, &rTex.vbo)
		gl.BindVertexArray(rTex.textureVao)

		gl.BindBuffer(gl.ARRAY_BUFFER, rTex.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, int(FLOAT_SIZE)*len(textureVertices), gl.Ptr(&textureVertices[0]), gl.STATIC_DRAW)

		gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 4*FLOAT_SIZE, 0)
		gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*FLOAT_SIZE, 2*FLOAT_SIZE)
		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)

		gl.GenVertexArrays(1, &rTex.outlineVao)
		gl.GenBuffers(1, &rTex.vbo)
		gl.BindVertexArray(rTex.outlineVao)

		gl.BindBuffer(gl.ARRAY_BUFFER, rTex.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, int(FLOAT_SIZE)*len(outlineVertices), gl.Ptr(&outlineVertices[0]), gl.STATIC_DRAW)

		gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 4*FLOAT_SIZE, 0)
		gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*FLOAT_SIZE, 2*FLOAT_SIZE)
		gl.EnableVertexAttribArray(0)
		gl.EnableVertexAttribArray(1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		fmt.Println(gl.CheckFramebufferStatus(gl.FRAMEBUFFER))
		panic("Framebuffer error")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
}

func CreateFramebuffer(width, height int32) (f *FrameBuffer) {
	f = new(FrameBuffer)
	f.width = width
	f.height = height
	gl.GenTextures(1, &f.colorBuffer)
	gl.BindTexture(gl.TEXTURE_2D, f.colorBuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, f.width, f.height, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.ActiveTexture(gl.TEXTURE1)

	gl.GenTextures(1, &f.depth)
	gl.BindTexture(gl.TEXTURE_2D, f.depth)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32, f.width, f.height, 0, gl.DEPTH_COMPONENT, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	gl.GenFramebuffers(1, &f.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fbo)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, f.colorBuffer, 0)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, f.depth, 0)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return f
}

func (f FrameBuffer) Render(clearColor [3]float32, objectColor [3]float32) {
	gl.Viewport(0, 0, f.width, f.height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, f.colorBuffer, 0)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, f.depth, 0)

	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(r.shaderHandle)
	gl.Uniform3f(r.uniforms.color, objectColor[0], objectColor[1], objectColor[2])
	gl.BindVertexArray(r.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (f *FrameBuffer) Resize(width, height int32) {
	f.width = width
	f.height = height
	gl.BindTexture(gl.TEXTURE_2D, f.colorBuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, f.width, f.height, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.BindTexture(gl.TEXTURE_2D, f.depth)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32, f.width, f.height, 0, gl.DEPTH_COMPONENT, gl.UNSIGNED_INT, nil)
}

func (f FrameBuffer) Image() uint32 {
	return f.colorBuffer
}

func (f FrameBuffer) Size() imgui.Vec2 {
	return imgui.NewVec2(float32(f.width), float32(f.height))
}

func Nuke() {
	gl.UseProgram(r.shaderHandle)
	gl.DeleteVertexArrays(1, &r.vao)
	gl.DeleteBuffers(1, &r.vbo)
	gl.DeleteProgram(r.shaderHandle)
}

func (f *FrameBuffer) RenderTexture(t1 uint32, zoom float32, pos [2]float32, width, height float32) {
	gl.Viewport(0, 0, f.width, f.height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fbo)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t1)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, f.colorBuffer, 0)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, f.depth, 0)

	size := f.Size()
	texAspect := width / height
	moveX := 2 * pos[0] / (texAspect * size.X * zoom)
	moveY := 2 * pos[1] / (size.Y * zoom)

	gl.UseProgram(rTex.shaderHandle)
	gl.Uniform1f(rTex.uniforms.aspect, float32(f.height)/float32(f.width)*texAspect)
	gl.Uniform1i(rTex.uniforms.outline, 0)
	gl.Uniform1f(rTex.uniforms.zoom, zoom)
	gl.Uniform2f(rTex.uniforms.move, moveX, moveY)
	gl.Uniform2f(rTex.uniforms.size, width, height)

	gl.ClearColor(0.29, 0.29, 0.39, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.BindVertexArray(rTex.textureVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.Uniform1i(rTex.uniforms.outline, 1)
	gl.BindVertexArray(rTex.outlineVao)
	gl.DrawArrays(gl.LINE_LOOP, 0, 4)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
