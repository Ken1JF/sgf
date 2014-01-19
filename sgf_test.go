package sgf

func ExampleSGFTypeSizes() {
    PrintSGFTypeSizes()
    // Output: 
    // Type Token size 1 alignment 1
    // Type ah.Position size 40 alignment 8
    // Type TreeNodeType size 1 alignment 1
    // Type TreeNodeIdx size 2 alignment 2
    // Type PropIdx size 2 alignment 2
    // Type TreeNode size 12 alignment 2
    // Type PropertyValue size 32 alignment 8
    // Type GameTree size 1520 alignment 8
    // Type Parser size 1848 alignment 8
    // Type PlayerInfo size 72 alignment 8
    // Type FF4Note size 1 alignment 1
    // Type SGFPropNodeType size 1 alignment 1
    // Type QualifierType size 1 alignment 1
    // Type PropValueType size 1 alignment 1
    // Type Property size 56 alignment 8
    // Type PropertyDefIdx size 1 alignment 1
    // Type ID_CountArray size 624 alignment 8
    // Type Scanner size 112 alignment 8
    // Type ErrorHandler size 8 alignment 8
    // Type ah.ErrorList size 24 alignment 8
    // Type Komi size 8 alignment 4
    // Type Result size 64 alignment 8
}
