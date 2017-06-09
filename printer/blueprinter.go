package printer

import (
	"blue/compiler"
	"blue/entities"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	Signal entities.Signal `json:"signal"`
	Index  int             `json:"index"`
}

type InstructionEntityGroup struct {
	InstructionProvider entities.ConstantCombinator
	InstructionDecider  entities.DeciderCombinator
	Lamp                entities.Lamp
}

var (
	signalInstructionIndex = entities.SignalItem("grenade")
	signalNegative         = entities.SignalItem("poison-capsule")
	signalAmux             = entities.SignalItem("crude-oil-barrel")
	signalBmux             = entities.SignalItem("lubricant-barrel")
	signalCmux             = entities.SignalItem("empty-barrel")
	signalCond             = entities.SignalItem("heavy-oil-barrel")
	signalAlu              = entities.SignalItem("light-oil-barrel")
	signalMbr              = entities.SignalItem("petroleum-gas-barrel")
	signalMar              = entities.SignalItem("water-barrel")
	signalRd               = entities.SignalItem("science-pack-1")
	signalWr               = entities.SignalItem("high-tech-science-pack")
	signalEnc              = entities.SignalItem("science-pack-2")
	signalA                = entities.SignalItem("science-pack-3")
	signalB                = entities.SignalItem("space-science-pack")
	signalC                = entities.SignalItem("production-science-pack")
	signalAddr             = entities.SignalItem("military-science-pack")
	signalImm              = entities.SignalItem("sulfuric-acid-barrel")
	signalBran             = entities.SignalItem("used-up-uranium-fuel-cell")
)

func (bp *BlueprintFile) addEntities(entities ...entities.IEntity) {
	bp.Blueprint.Entities = append(bp.Blueprint.Entities, entities...)
}

func CreateBlueprint(instructions []compiler.Instruction, writer io.Writer) error {

	//Build Bp file
	bpf := BlueprintFile{
		Blueprint: Blueprint{
			Icons: []Icon{
				Icon{Signal: entities.SignalItem("science-pack-1"), Index: 1},
				Icon{Signal: entities.SignalItem("science-pack-2"), Index: 2},
				Icon{Signal: entities.SignalItem("science-pack-3"), Index: 3},
				Icon{Signal: entities.SignalItem("space-science-pack"), Index: 4},
			},
			Entities: []entities.IEntity{},
			Item:     "blueprint",
			Version:  64425558017,
		},
	}

	//Create init stuff
	addBusConnector(&bpf)

	//Build and connect each row
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

	//Create string
	str, err := bpf.asBlueprintString()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(str))
	if err != nil {
		return err
	}

	return nil
}

func addBusConnector(bp *BlueprintFile) {
	instructionIndexInput := entities.NewArithmaticCombinatorWithConstant(signalInstructionIndex, 0, entities.OperationAddition, signalInstructionIndex)
	instructionIndexInput.Position.Set(0, 1.5)

	allReturn := entities.NewArithmaticCombinatorWithConstant(entities.SignalVirtualEach, 0, entities.OperationAddition, entities.SignalVirtualEach)
	allReturn.Position.Set(1, 1.5)
	allReturn.Direction = entities.DirectionDown

	negateInstructionIndexReturn := entities.NewArithmaticCombinator(signalInstructionIndex, signalNegative, entities.OperationMultiplication, signalInstructionIndex)
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
			signalBran, instruction.Bran,
		),
		InstructionDecider: entities.NewDeciderCombinatorWithConstant(signalInstructionIndex, index, entities.ComparatorEqual, entities.SignalVirtualEverything, true),
		Lamp:               entities.NewLampWithConstant(entities.SignalVirtualAnything, 0, entities.ComparatorGreaterThan),
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

func (bp BlueprintFile) asBlueprintString() (outString string, err error) {
	var jsonBuffer bytes.Buffer
	jsonEncoder := json.NewEncoder(&jsonBuffer)
	jsonEncoder.SetEscapeHTML(false)

	err = jsonEncoder.Encode(bp)
	if err != nil {
		return
	}

	fmt.Println(jsonBuffer.String())

	var zippedBuffer bytes.Buffer
	gzipper := zlib.NewWriter(&zippedBuffer)

	_, err = gzipper.Write(jsonBuffer.Bytes())
	gzipper.Close()
	if err != nil {
		return
	}

	outString = "0" + base64.StdEncoding.EncodeToString(zippedBuffer.Bytes())
	return
}
