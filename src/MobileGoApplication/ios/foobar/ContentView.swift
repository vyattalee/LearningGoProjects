// ios/foobar/ContentView.swift

struct ContentView: View {

    @State private var txt: String = ""

    var body: some View {
        VStack{
            TextField("", text: $txt)
            .textFieldStyle(RoundedBorderTextFieldStyle())
            Button("Reverse"){
                // Reverse text here
            }
            Spacer()
        }
        .padding(.all, 15)
    }
}

// ios/foobar/ContentView.swift

Button("Reverse"){
    let str = reverse(UnsafeMutablePointer<Int8>(mutating: (self.txt as NSString).utf8String))
    self.txt = String.init(cString: str!, encoding: .utf8)!
    // don't forget to release the memory to the C String
    str?.deallocate()
}