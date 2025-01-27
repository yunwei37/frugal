/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ssa

import (
    `github.com/oleiade/lane`
)

type CFG struct {
    Root              *BasicBlock
    Depth             map[int]int
    DominatedBy       map[int]*BasicBlock
    DominatorOf       map[int][]*BasicBlock
    DominanceFrontier map[int][]*BasicBlock
}

func (self *CFG) Rebuild() {
    updateDominatorTree(self)
    updateDominatorDepth(self)
    updateDominatorFrontier(self)
}

func (self *CFG) MaxBlock() int {
    var id int
    var ret int

    /* get the max ID */
    for id = range self.DominatedBy {
        if id > ret {
            ret = id
        }
    }

    /* select between ID and root ID */
    if ret > self.Root.Id {
        return ret
    } else {
        return self.Root.Id
    }
}

func (self *CFG) PostOrder(action func(bb *BasicBlock)) {
    stack := lane.NewStack()
    visited := make(map[int]bool)

    /* add root node */
    visited[self.Root.Id] = true
    stack.Push(self.Root)

    /* traverse the graph */
    for !stack.Empty() {
        tail := true
        this := stack.Head().(*BasicBlock)

        /* add all the successors */
        for _, p := range self.DominatorOf[this.Id] {
            if !visited[p.Id] {
                tail = false
                visited[p.Id] = true
                stack.Push(p)
                break
            }
        }

        /* all the successors are visited, pop the current node */
        if tail {
            action(stack.Pop().(*BasicBlock))
        }
    }
}

func (self *CFG) ReversePostOrder(action func(bb *BasicBlock)) {
    var i int
    var bb []*BasicBlock

    /* traverse as post-order */
    self.PostOrder(func(p *BasicBlock) {
        bb = append(bb, p)
    })

    /* reverse post-order */
    for i = len(bb) - 1; i >= 0; i-- {
        action(bb[i])
    }
}
