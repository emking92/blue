package blueprinter

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"factorio-assembly/compiler"
	"factorio-assembly/entities"
	"factorio-assembly/entities/controls"
	"fmt"
)

type BlueprintFile struct {
	Blueprint Blueprint `json:"blueprint"`
}
type Blueprint struct {
	Icons                        []Icon             `json:"icons"`
	Entities                     []entities.IEntity `json:"entities"`
	Item                         string             `json:"item"`
	Version                      int64              `json:"version"`
	instructionIndexInput        *entities.ArithmaticCombinator
	allReturn                    *entities.ArithmaticCombinator
	negateInstructionIndexReturn *entities.ArithmaticCombinator
	negativeOne                  *entities.ConstantCombinator
}

type Icon struct {
	Signal controls.Signal `json:"signal"`
	Index  int             `json:"index"`
}

type InstructionEntityGroup struct {
	InstructionProvider entities.ConstantCombinator
	InstructionDecider  entities.DeciderCombinator
	Lamp                entities.Lamp
}

var (
	signalInstructionIndex = controls.SignalItem("grenade")
	signalNegative         = controls.SignalItem("poison-capsule")
	signalAmux             = controls.SignalItem("crude-oil-barrel")
	signalBmux             = controls.SignalItem("lubricant-barrel")
	signalCmux             = controls.SignalItem("empty-barrel")
	signalCond             = controls.SignalItem("heavy-oil-barrel")
	signalAlu              = controls.SignalItem("light-oil-barrel")
	signalMbr              = controls.SignalItem("petroleum-gas-barrel")
	signalMar              = controls.SignalItem("water-barrel")
	signalRd               = controls.SignalItem("science-pack-1")
	signalWr               = controls.SignalItem("high-tech-science-pack")
	signalEnc              = controls.SignalItem("science-pack-2")
	signalA                = controls.SignalItem("science-pack-3")
	signalB                = controls.SignalItem("space-science-pack")
	signalC                = controls.SignalItem("production-science-pack")
	signalAddr             = controls.SignalItem("military-science-pack")
	signalImm              = controls.SignalItem("sulfuric-acid-barrel")
)

func (bp *BlueprintFile) addEntities(entities ...entities.IEntity) {
	bp.Blueprint.Entities = append(bp.Blueprint.Entities, entities...)
}

func CreateBlueprint( /*instructions []compiler.Instruction, writer io.Writer*/ ) {

	//Build Bp file
	bpf := BlueprintFile{
		Blueprint: Blueprint{
			Icons: []Icon{
				Icon{Signal: controls.SignalItem("science-pack-1"), Index: 1},
				Icon{Signal: controls.SignalItem("science-pack-2"), Index: 2},
				Icon{Signal: controls.SignalItem("science-pack-3"), Index: 3},
				Icon{Signal: controls.SignalItem("space-science-pack"), Index: 4},
			},
			Entities: []entities.IEntity{},
			Item:     "blueprint",
			Version:  64425558017,
		},
	}

	//Create init stuff
	addBusConnector(&bpf)

	//Build and connect each row
	instructions := []compiler.Instruction{
		compiler.Instruction{C: 11, Enc: 1, Alu: 8},
		compiler.Instruction{C: 10, Enc: 1, Cmux: 1, Imm: 13},
		compiler.Instruction{C: 10, A: 10, B: 8, Enc: 1, Alu: 0},
		compiler.Instruction{C: 11, Enc: 1, Cmux: 1, Imm: 2},
		compiler.Instruction{C: 11, Enc: 1, A: 11, B: 10, Alu: 2},
		compiler.Instruction{Cond: 3},
	}

	instructionEntityGroups := make([]InstructionEntityGroup, len(instructions))

	for i, instruction := range instructions {
		group := buildEntitiesFromInstruction(&bpf, instruction, i)
		instructionEntityGroups[i] = group

		if i == 0 {
			entities.ConnectEntities(group.InstructionDecider.Entity, 1, bpf.Blueprint.instructionIndexInput.Entity, 2, entities.WireGreen)
			entities.ConnectEntities(group.InstructionDecider.Entity, 2, bpf.Blueprint.negateInstructionIndexReturn.Entity, 1, entities.WireGreen)
		} else {
			entities.ConnectEntities(group.InstructionDecider.Entity, 1, instructionEntityGroups[i-1].InstructionDecider.Entity, 1, entities.WireGreen)
			entities.ConnectEntities(group.InstructionDecider.Entity, 2, instructionEntityGroups[i-1].InstructionDecider.Entity, 2, entities.WireGreen)
		}
	}

	//Create file
	bpf.asBlueprintString()

}

func addBusConnector(bp *BlueprintFile) {
	instructionIndexInput := entities.NewArithmaticCombinatorWithConstant(signalInstructionIndex, 0, controls.OperationAddition, signalInstructionIndex)
	instructionIndexInput.Position.Set(0, 1.5)

	allReturn := entities.NewArithmaticCombinatorWithConstant(controls.SignalVirtualEach, 0, controls.OperationAddition, controls.SignalVirtualEach)
	allReturn.Position.Set(1, 1.5)
	allReturn.Direction = entities.DirectionDown

	negateInstructionIndexReturn := entities.NewArithmaticCombinator(signalInstructionIndex, signalNegative, controls.OperationMultiplication, signalInstructionIndex)
	negateInstructionIndexReturn.Position.Set(2, 1.5)
	negateInstructionIndexReturn.Direction = entities.DirectionDown

	negativeOne := entities.NewConstantCombinator(signalNegative, -1)
	negativeOne.Position.Set(3, 1)

	entities.ConnectEntities(instructionIndexInput.Entity, 1, allReturn.Entity, 2, entities.WireGreen)
	entities.ConnectEntities(allReturn.Entity, 2, negateInstructionIndexReturn.Entity, 2, entities.WireGreen)
	entities.ConnectEntities(allReturn.Entity, 1, negateInstructionIndexReturn.Entity, 1, entities.WireGreen)
	entities.ConnectEntities(negativeOne.Entity, 1, negateInstructionIndexReturn.Entity, 1, entities.WireRed)

	bp.addEntities(&instructionIndexInput, &allReturn, &negateInstructionIndexReturn, &negativeOne)
	bp.Blueprint.allReturn = &allReturn
	bp.Blueprint.instructionIndexInput = &instructionIndexInput
	bp.Blueprint.negateInstructionIndexReturn = &negateInstructionIndexReturn
	bp.Blueprint.negativeOne = &negativeOne
}

func buildEntitiesFromInstruction(bp *BlueprintFile, instruction compiler.Instruction, index int) InstructionEntityGroup {
	entityGroup := InstructionEntityGroup{
		InstructionProvider: entities.NewConstantCombinator(
			signalAmux, instruction.Amux,
			signalCmux, instruction.Cmux,
			signalCond, instruction.Cond,
			signalAlu, instruction.Alu,
			signalBmux, instruction.Bmux,
			signalMbr, instruction.Mbr,
			signalMar, instruction.Mar,
			signalRd, instruction.Rd,
			signalWr, instruction.Wr,
			signalEnc, instruction.Enc,
			signalA, instruction.A,
			signalB, instruction.B,
			signalC, instruction.C,
			signalAddr, instruction.Addr,
			signalImm, instruction.Imm,
		),
		InstructionDecider: entities.NewDeciderCombinatorWithConstant(signalInstructionIndex, index, controls.ComparatorEqual, controls.SignalVirtualEverything, true),
		Lamp:               entities.NewLampWithConstant(controls.SignalVirtualAnything, 0, controls.ComparatorGreaterThan),
	}

	entityGroup.InstructionProvider.Direction = entities.DirectionRight
	entityGroup.InstructionProvider.Position.Set(0, float32(-index))

	entityGroup.InstructionDecider.Direction = entities.DirectionRight
	entityGroup.InstructionDecider.Position.Set(1.5, float32(-index))

	entityGroup.Lamp.Position.Set(3, float32(-index))

	entities.ConnectEntities(entityGroup.InstructionProvider.Entity, 1, entityGroup.InstructionDecider.Entity, 1, entities.WireRed)
	entities.ConnectEntities(entityGroup.InstructionDecider.Entity, 2, entityGroup.Lamp.Entity, 1, entities.WireRed)

	bp.addEntities(&entityGroup.InstructionProvider, &entityGroup.InstructionDecider, &entityGroup.Lamp)
	return entityGroup
}

func (bp BlueprintFile) asBlueprintString() {
	var jsonBuffer bytes.Buffer
	jsonEncoder := json.NewEncoder(&jsonBuffer)
	jsonEncoder.SetEscapeHTML(false)

	err := jsonEncoder.Encode(bp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(jsonBuffer.String())

	var zippedBuffer bytes.Buffer
	gzipper := zlib.NewWriter(&zippedBuffer)
	_, err = gzipper.Write(jsonBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	gzipper.Close()

	outString := base64.StdEncoding.EncodeToString(zippedBuffer.Bytes())

	fmt.Println(outString)
}
