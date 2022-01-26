package tests

import (
    `bytes`
    `fmt`
    `os/exec`
    `plugin`
    `reflect`
    `runtime`
    `runtime/debug`
    `sync`
    `testing`

    _ `github.com/cloudwego/frugal`
    `github.com/cloudwego/frugal/testdata/kitex_gen/baseline`
)

func pluginInit() {
    bin, err := exec.LookPath("go")
    if err != nil {
        panic(err)
    }
    out := bytes.NewBuffer(nil)
    cmd := exec.Cmd{
        Path: bin,
        Args: []string{"go", "build", "-buildmode", "plugin", "-o", "plugin/plugin."+runtime.Version()+".so", "plugin/main.go"},
        Stdout: out,
        Stderr: out,
    }
    if err := cmd.Run(); err != nil {
        panic(string(out.Bytes()))
    }
}

func pluginTestMain() {
    go func ()  {
        println("Begin GC looping...")
        for {
            runtime.GC()
            debug.FreeOSMemory() 
        }
	}()
	runtime.Gosched()
}

func TestPlugin(t *testing.T) {
    pluginInit()
    pluginTestMain()
    p, err := plugin.Open("plugin/plugin."+runtime.Version()+".so")
    if err != nil {
        t.Fatal(err)
    }
    v, err := p.Lookup("V")
    if err != nil {
        t.Fatal(err)
    }
    f, err := p.Lookup("F")
    if err != nil {
        t.Fatal(err)
    }
    *v.(*int) = 7
    f.(func())() // prints "Hello, number 7"
    obj, err := p.Lookup("Obj")
    m := *(obj.(*[]byte))
    fmt.Printf("%#v\n", m)

    wg := sync.WaitGroup{}
    for i:=0; i<1000; i++ {
        wg.Add(1)
        go func(){
            defer wg.Done()
            d, err := p.Lookup("Marshal")
            if err != nil {
                t.Error(err)
                return
            }
            enc := d.(func(val interface{}) ([]byte, error))
            var exp baseline.Simple
            out, err := enc(&exp)
            if err != nil {
                t.Error(err)
                return
            }
            if !reflect.DeepEqual(m, out) {
                t.Error(m, out)
                return
            }
        }()
    }
    wg.Wait()
}