package animation

import (
	"VulpesEditor/app/front/renderer"
	"VulpesEditor/app/objectModel/model"
	"fmt"
	"math"
	"slices"
	"time"

	im "github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl32"
)

type AnimationContext struct {
	model *model.Model
	mesh  *renderer.Mesh

	bones     []*bone
	animation *animation
	boneAnim  *boneAnimation

	animator *animator

	running bool

	editBone          *bone
	editKeyframe      *keyFrame
	bonesToSelect     []string
	keyFramesToSelect []string
	boneSelected      int32
	keyFrameSelected  int32
}

type bone struct {
	unitId   string
	boneId   uint32
	parent   *bone
	children []*bone
}

func newBone(unit *model.CubeUnit) (b *bone) {
	b = new(bone)
	b.unitId = unit.Id
	b.boneId = unit.BoneID
	return
}

func calcBones() {
	ctx.bones = nil
	ctx.bonesToSelect = nil
	for _, u := range ctx.model.Units {
		ctx.bones = append(ctx.bones, newBone(u))
		ctx.bonesToSelect = append(ctx.bonesToSelect, u.Name)
		bA := new(boneAnimation)
		bA.unitId = u.Id
		ctx.animation.bones[u.Id] = bA
	}
	for _, u := range ctx.model.Units {
		var current *bone
		for _, b := range ctx.bones {
			if b.unitId == u.Id {
				current = b
				break
			}
		}
		parent, ok := u.Parent().(*model.CubeUnit)
		if ok {
			for _, b := range ctx.bones {
				if b.unitId == parent.Id {
					current.parent = b
					b.children = append(b.children, current)
					break
				}
			}
		}
	}
}

type keyFrame struct {
	id     uint32
	unitId string
	time   float32
	trans  transformation
}

func newFrame(id uint32) (k *keyFrame) {
	k = new(keyFrame)
	k.id = id
	return
}

type animation struct {
	id    string
	name  string
	bones map[string]*boneAnimation
}

type boneAnimation struct {
	unitId    string
	keyFrames []*keyFrame
	count     uint32
}

func (s *boneAnimation) addKeyFrame() (f *keyFrame) {
	s.count++
	f = newFrame(s.count)
	f.unitId = s.unitId
	f.trans.size[0] = 1
	f.trans.size[1] = 1
	f.trans.size[2] = 1
	s.keyFrames = append(s.keyFrames, f)
	return
}

type transformation struct {
	pos  [3]float32
	size [3]float32
	rot  [3]float32
}

func selectBone(index int32) {
	if int(index) >= len(ctx.bones) {
		return
	}

	ctx.editBone = ctx.bones[index]
	var ok bool
	var unitId = ctx.editBone.unitId
	ctx.boneAnim, ok = ctx.animation.bones[unitId]
	if !ok {
		ctx.boneAnim = new(boneAnimation)
		ctx.boneAnim.unitId = unitId
		ctx.animation.bones[unitId] = ctx.boneAnim
	}
}

func selectKeyFrame(index int32) {
	if int(index) >= len(ctx.boneAnim.keyFrames) {
		return
	}
	ctx.editKeyframe = ctx.boneAnim.keyFrames[index]
}

func setKeyFrameToEdit(k *keyFrame) {
	if ctx.boneAnim != nil {
		for idx, k2 := range ctx.boneAnim.keyFrames {
			if k2 == k {
				selectKeyFrame(int32(idx))
			}
		}
	}
	ctx.editKeyframe = k
}

type boneTransformation struct {
	time int64

	trans transformation
}

type boneAnimator struct {
	keyframes []boneTransformation
}

func linear(n1, n2, t float32) float32 {
	return (n1-n2)*t + n2
}

func (s *boneAnimator) get(time int64) (ok bool, trans mgl32.Mat4) {
	if len(s.keyframes) == 0 {
		ok = false
		return
	}
	idx := 0
	last := s.keyframes[0]
	var next boneTransformation
	for {
		idx++
		if idx >= len(s.keyframes) {
			next = s.keyframes[0]
			return
		}
		next = s.keyframes[idx]
		if next.time > time {
			break
		}
		last = next
	}
	var b transformation
	initTime := last.time
	endTime := next.time
	coeff := float32(time-initTime) / float32(endTime-initTime)

	b.pos[0] = linear(last.trans.pos[0], next.trans.pos[0], coeff)
	b.pos[1] = linear(last.trans.pos[1], next.trans.pos[1], coeff)
	b.pos[2] = linear(last.trans.pos[2], next.trans.pos[2], coeff)

	b.size[0] = linear(last.trans.size[0], next.trans.size[0], coeff)
	b.size[1] = linear(last.trans.size[1], next.trans.size[1], coeff)
	b.size[2] = linear(last.trans.size[2], next.trans.size[2], coeff)

	b.rot[0] = linear(last.trans.rot[0], next.trans.rot[0], coeff)
	b.rot[1] = linear(last.trans.rot[1], next.trans.rot[1], coeff)
	b.rot[2] = linear(last.trans.rot[2], next.trans.rot[2], coeff)

	trans = mgl32.Ident4()
	trans = trans.Mul4(mgl32.Scale3D(b.size[0], b.size[1], b.size[2]))
	trans = trans.Mul4(mgl32.HomogRotate3D(b.rot[0], mgl32.Vec3{1, 0, 0}))
	trans = trans.Mul4(mgl32.HomogRotate3D(b.rot[1], mgl32.Vec3{0, 1, 0}))
	trans = trans.Mul4(mgl32.HomogRotate3D(b.rot[2], mgl32.Vec3{0, 0, 1}))
	trans = trans.Mul4(mgl32.Translate3D(b.pos[0], b.pos[1], b.pos[2]))

	ok = true
	return
}

type animator struct {
	loopTime int64
	time     int64
	step     float32

	bones map[uint32]boneAnimator

	last int64
}

func (s *animator) loop() {
	now := time.Now().UnixMilli()
	s.time += now - s.last
	s.last = now

	if s.time > s.loopTime {
		s.time = s.time % s.loopTime
	}
}

func (s *animator) getTrans(boneId uint32) mgl32.Mat4 {
	trans, ok1 := s.bones[boneId]
	if !ok1 {
		return mgl32.Ident4()
	}
	ok2, m := trans.get(s.time)
	if !ok2 {
		return mgl32.Ident4()
	}
	return m
}

func newAnimator() (r *animator) {
	r = new(animator)
	r.last = time.Now().UnixMilli()
	r.bones = make(map[uint32]boneAnimator)
	return
}

func constructAnimation() {
	ctx.animator = newAnimator()
	var maximum int64
	for _, b := range ctx.bones {
		var bA boneAnimator
		boneA, ok := ctx.animation.bones[b.unitId]
		if !ok {
			continue
		}
		for _, k := range boneA.keyFrames {
			if k.unitId == b.unitId {
				var key boneTransformation
				key.time = int64(k.time * 1000)
				maximum = max(maximum, key.time)
				key.trans = k.trans
				bA.keyframes = append(bA.keyframes, key)
			}
		}
		if len(bA.keyframes) == 0 {
			continue
		}
		slices.SortFunc(bA.keyframes, func(A, B boneTransformation) int {
			return int(A.time - B.time)
		})
		ctx.animator.bones[b.boneId] = bA
	}
	ctx.animator.loopTime = maximum
	if ctx.animator.loopTime == 0 {
		ctx.animator = nil
	}
}

func Show(id string) {
	ctxManager.Check(id)

	ctx.keyFramesToSelect = nil
	if ctx.boneAnim != nil {
		for _, k := range ctx.boneAnim.keyFrames {
			ctx.keyFramesToSelect = append(ctx.keyFramesToSelect, fmt.Sprintf("Keyframe #%d", k.id))
		}
	}

	if ctx.running {
		if ctx.animator != nil {
			ctx.animator.loop()
			for _, b := range ctx.bones {
				m := ctx.animator.getTrans(b.boneId)
				ctx.mesh.SetBoneTransMatrix(b.boneId, m)
			}
		}
	}

	if ctx.model.Changed {
		calcBones()
		selectBone(0)
		f1 := ctx.boneAnim.addKeyFrame()
		f1.time = 0
		f1.trans.rot[2] = 0
		f2 := ctx.boneAnim.addKeyFrame()
		f2.time = 2
		f2.trans.rot[2] = 2 * math.Pi
		// f3 := ctx.boneAnim.addKeyFrame()
		// f3.time = 1
		// f3.trans.rot[2] = 0
		// f4 := ctx.boneAnim.addKeyFrame()
		// f4.time = 1.5
		// f4.trans.rot[2] = 3 * math.Pi / 2
		// f5 := ctx.boneAnim.addKeyFrame()
		// f5.time = 2
		// f5.trans.rot[2] = 2 * math.Pi

		constructAnimation()

		ctx.running = false
	}

	if im.Begin("Animation") {
		if im.Button("Run") {
			ctx.running = !ctx.running
		}
		if im.ComboStrarr("Bone", &ctx.boneSelected, ctx.bonesToSelect, int32(len(ctx.bonesToSelect))) {
			selectBone(ctx.boneSelected)
		}
		cant := ctx.boneAnim == nil
		if cant {
			im.BeginDisabled()
		}
		if im.Button("Add Keyframe") {
			setKeyFrameToEdit(ctx.boneAnim.addKeyFrame())
		}
		if im.ComboStrarr("Keyframe", &ctx.keyFrameSelected, ctx.keyFramesToSelect, int32(len(ctx.keyFramesToSelect))) {
			selectKeyFrame(ctx.keyFrameSelected)
		}
		if cant {
			im.EndDisabled()
		}
		if ctx.editKeyframe != nil {
			c1 := im.InputFloat3("Position", &ctx.editKeyframe.trans.pos)
			c2 := im.InputFloat3("Size", &ctx.editKeyframe.trans.size)
			c3 := im.InputFloat3("Rotation", &ctx.editKeyframe.trans.rot)
			c4 := im.InputFloat("Time", &ctx.editKeyframe.time)
			if c1 || c2 || c3 || c4 {
				constructAnimation()
			}
		}
	}
	im.End()
}
