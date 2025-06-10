<template>
  <div class="p-2 m-2 flex gap-4">
    <button
      class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
      @click="addRowToCurrent"
    >
      Add Row
    </button>
    <button
      class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
      @click="addColumnToCurrent"
    >
      Add Column
    </button>
  </div>
  <div
    @scroll="handleScroll"
    ref="scrollContainer"
    class="overflow-auto relative h-96"
  >
    <table
      ref="table"
      class="excel-grid border-collapse table-fixed border border-gray-300 select-none h-fit"
      contenteditable="false"
      @copy="handleCopy"
      @cut="handleCut"
      @mousemove="handleMouseMove"
      @mouseup="endSelection"
    >
      <thead>
        <tr>
          <th></th>
          <th
            class="border border-gray-300 bg-gray-100 p-2 text-left"
            v-for="(column, index) in columns"
            :key="index"
          >
            <div class="flex justify-between items-center w-full gap-2">
              <span
                v-if="!isEditingColumn(index)"
                @click="editColumnName(index)"
                class="cursor-pointer"
              >
                {{ column }}
              </span>
              <input
                v-else
                type="text"
                v-model="columns[index]"
                @blur="stopEditingColumn"
                @keydown.enter="stopEditingColumn"
                class="border border-gray-400 p-1 w-full text-sm"
              />
              <input
                type="text"
                class="border border-gray-400 p-1 w-full text-sm"
                placeholder="Filter"
                v-model="filters[index]"
                @input="applyFilters"
              />
              <span @click="toggleSort(index)" class="cursor-pointer">
                {{ sortOrder === "asc" ? "▲" : "▼" }}
              </span>
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(row, rowIndex) in visibleRows"
          :key="rowIndex + startRow"
          class="outline-hidden"
          contenteditable="true"
        >
          <th class="border border-gray-300 bg-gray-100 p-2">
            {{ rowIndex + 1 + startRow }}
          </th>
          <td
            :class="{
              'bg-blue-100': isSelectedCell(rowIndex, colIndex),
            }"
            class="border border-gray-300 p-2"
            v-for="(value, colIndex) in row"
            :key="colIndex"
            @paste="handlePaste"
            @mousemove="dragCell(rowIndex, colIndex)"
            @mousedown="startCell(rowIndex, colIndex)"
            @mouseup="handleCellClick(rowIndex, colIndex, $event)"
          >
            <span v-if="!isEditingCell(rowIndex, colIndex)">
              {{ value }}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <input
    type="text"
    ref="dynamicInput"
    @paste="handlePaste"
    v-model="currentInputValue"
    class="absolute bg-white border border-blue-500 p-1 text-sm"
    v-if="isInputVisible"
    :style="inputStyle"
    @blur="hideInput"
    @keydown="handleKeyDown"
  />
</template>
<script setup lang="ts">
import { computed, nextTick, reactive, ref } from "vue"
import { ClipboardEvent } from "happy-dom"

const dummyColumns = ["ID", "Name", "Age", "Country"]
const dummyData1 = Array.from({ length: 1000000 }, (_, i) => [
  i + 1,
  `User ${i + 1}`,
  20 + (i % 30),
  `Country ${i % 5}`,
])
const columns = ref(dummyColumns)
const rowData = ref(dummyData1)

// 필터링 및 정렬된 데이터
const filters = reactive(Array(columns.value.length).fill(""))

// 캐시된 정렬 결과
const sortedCache = ref<Record<string, never[]>>({})

const applySort = (data: never) => {
  if (sortColumn.value === null || sortOrder.value === null) {
    return data
  }

  const cacheKey = `${sortColumn.value}-${sortOrder.value}`

  // 캐시 확인
  if (sortedCache.value[cacheKey]) {
    return sortedCache.value[cacheKey]
  }

  // 정렬 수행
  const sortedData = [...data].sort((a, b) => {
    const valA = a[sortColumn.value!]
    const valB = b[sortColumn.value!]

    if (valA < valB) return sortOrder.value === "asc" ? -1 : 1
    if (valA > valB) return sortOrder.value === "asc" ? 1 : -1
    return 0
  })

  // 캐시 저장
  sortedCache.value[cacheKey] = sortedData
  return sortedData
}

const filteredData = computed(() => {
  let data = rowData.value

  // 필터 적용
  filters.forEach((filter, index) => {
    if (filter.trim()) {
      data = data.filter((row) =>
        String(row[index]).toLowerCase().includes(filter.toLowerCase())
      )
    }
  })

  // 정렬 적용
  // @ts-ignore
  data = applySort(data)

  return data
})

// 정렬 토글
const toggleSort = (columnIndex: number) => {
  if (sortColumn.value === columnIndex) {
    sortOrder.value =
      sortOrder.value === "asc"
        ? "desc"
        : sortOrder.value === "desc"
          ? null
          : "asc"
  } else {
    sortColumn.value = columnIndex
    sortOrder.value = "asc"
  }
}

const rowHeight = 40 // 행의 높이
const containerHeight = 400 // 가시 영역 높이

const scrollTop = ref(0) // 스크롤 위치

const visibleCount = computed(() => Math.ceil(containerHeight / rowHeight))

const startRow = computed(() => Math.floor(scrollTop.value / rowHeight))

const endRow = computed(() =>
  Math.min(filteredData.value.length, startRow.value + visibleCount.value + 2)
)

const sortColumn = ref<number | null>(null)
const sortOrder = ref<"asc" | "desc" | null>(null)

const visibleRows = computed(() =>
  filteredData.value.slice(startRow.value, endRow.value)
)

const handleScroll = (event: Event) => {
  scrollTop.value = (event.target as HTMLElement)?.scrollTop
}

const applyFilters = () => {
  scrollTop.value = 0 // 필터 적용 시 스크롤 위치 초기화
}

const currentEditingColumnIndex = ref<number | null>(null) // 현재 편집 중인 열 인덱스
const editColumnName = (index: number) => {
  currentEditingColumnIndex.value = index // 열 편집 시작
}

const stopEditingColumn = () => {
  currentEditingColumnIndex.value = null // 열 편집 종료
}

const isEditingColumn = (index: number) => {
  return currentEditingColumnIndex.value === index
}

const currentInputValue = ref<string | number>("")
const currentEditingCell = ref<Position | null>(null) // 현재 편집 중인 셀
const isInputVisible = ref(false)
const inputStyle = reactive({
  top: "0px",
  left: "0px",
  width: "0px",
  height: "0px",
  zIndex: -1,
})

const table = ref<HTMLTableElement | null>(null)
const dynamicInput = ref<HTMLInputElement | null>(null)

const showInput = (rowIndex: number, colIndex: number) => {
  if (hasDragged.value) return // 드래그 중에는 입력창 표시하지 않음

  selectedCells.splice(0)
  selectedCells.push({ rowIndex, colIndex })
  currentEditingCell.value = { rowIndex, colIndex }
  currentCellPosition.value = { rowIndex, colIndex }
  currentInputValue.value = rowData.value[rowIndex][colIndex]

  const cell = table.value!.rows[rowIndex + 1]?.cells[colIndex + 1] // +1 for header row
  if (!cell) return

  const cellRect = cell.getBoundingClientRect()
  inputStyle.top = `${cellRect.top + window.scrollY}px`
  inputStyle.left = `${cellRect.left + window.scrollX}px`
  inputStyle.width = `${cellRect.width}px`
  inputStyle.height = `${cellRect.height}px`
  inputStyle.zIndex = 10
  isInputVisible.value = true

  nextTick(() => {
    dynamicInput.value?.focus()
  })
}

const handleKeyDown = (event: Event) => {
  const keyboardEvent = event as KeyboardEvent
  switch (keyboardEvent.key) {
    case "Tab":
      event.preventDefault()
      // Shift + Tab: Move to the previous cell
      navigateInput(keyboardEvent.shiftKey ? -1 : 1)
      break
    case "Enter":
      saveInput()
      break
    case "Escape":
      hideInput()
      break
  }
}

const navigateInput = (direction: number) => {
  if (!currentEditingCell.value) return

  const { rowIndex, colIndex } = currentEditingCell.value
  const totalCols = columns.value.length
  const totalRows = rowData.value.length

  let newRowIndex = rowIndex
  let newColIndex = colIndex + direction

  if (newColIndex < 0) {
    newRowIndex--
    newColIndex = totalCols - 1
  } else if (newColIndex >= totalCols) {
    newRowIndex++
    newColIndex = 0
  }

  if (newRowIndex >= 0 && newRowIndex < totalRows) {
    showInput(newRowIndex, newColIndex)
  }
}

const hasDragged = ref(false)
const handleCellClick = (
  rowIndex: number,
  colIndex: number,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  event: Event
) => {
  if (!hasDragged.value) {
    // 드래그가 아닌 클릭일 때만 showInput 호출
    showInput(rowIndex, colIndex)
  }
}

const hideInput = () => {
  isInputVisible.value = false
  currentEditingCell.value = null
  inputStyle.zIndex = -1
}

const saveInput = () => {
  const { rowIndex, colIndex } = currentEditingCell.value!
  rowData.value[rowIndex][colIndex] = currentInputValue.value
  hideInput()
}

const isEditingCell = (rowIndex: number, colIndex: number) => {
  return (
    currentEditingCell.value &&
    currentEditingCell.value.rowIndex === rowIndex &&
    currentEditingCell.value.colIndex === colIndex
  )
}

const currentCellPosition = ref<Position | null>(null)
const currentRowIndex = computed(
  () => currentCellPosition.value?.rowIndex ?? -1
)
const currentColIndex = computed(
  () => currentCellPosition.value?.colIndex ?? -1
)

const addRow = (row = 0) => {
  // Adds a new empty row
  const newRow = Array(columns.value.length).fill("")
  rowData.value.splice(row, 0, newRow)
}

const addColumn = (col = 0) => {
  // Adds a new column with default header
  const newColumnName = `Column ${columns.value.length + 1}`
  columns.value.splice(col, 0, newColumnName)
  rowData.value.forEach((row) => row.splice(col, 0, ""))
}
const addColumnToCurrent = () => {
  addColumn(currentColIndex.value + 1)
}

const addRowToCurrent = () => {
  addRow(currentRowIndex.value + 1)
}

const handleCut = (event: Event) => {
  // 잘라내기 이벤트 처리
  event.preventDefault()

  // 복사 이벤트와 동일한 로직
  handleCopy(event)

  // 선택된 셀 초기화
  selectedCells.forEach((cell) => {
    rowData.value[cell.rowIndex][cell.colIndex] = ""
  })
}

const handleCopy = (event: Event) => {
  event.preventDefault()

  // Find the range of selected rows and columns
  const minRow = Math.min(...selectedCells.map((cell) => cell.rowIndex))
  const maxRow = Math.max(...selectedCells.map((cell) => cell.rowIndex))
  const minCol = Math.min(...selectedCells.map((cell) => cell.colIndex))
  const maxCol = Math.max(...selectedCells.map((cell) => cell.colIndex))

  // Generate text to copy
  let copyText = ""
  for (let r = minRow; r <= maxRow; r++) {
    let rowText = []
    for (let c = minCol; c <= maxCol; c++) {
      if (
        selectedCells.some((cell) => cell.rowIndex === r && cell.colIndex === c)
      ) {
        rowText.push(rowData.value[r][c])
      } else {
        rowText.push("") // Empty for non-selected cells in range
      }
    }
    copyText += rowText.join("\t") + "\n"
  }

  ;(event as unknown as ClipboardEvent).clipboardData?.setData(
    "text/plain",
    copyText.trim()
  )
}

const handlePaste = (event: Event) => {
  // 붙여넣기 이벤트 처리
  event.preventDefault()
  const clipboardEvent = event as unknown as ClipboardEvent
  const pastedData = clipboardEvent.clipboardData?.getData("text/plain") || ""
  const rows = pastedData.split("\n").map((row: string) => row.split("\t"))

  const { rowIndex, colIndex } = currentEditingCell.value!

  rows.forEach((row, rowOffset) => {
    row.forEach((value, colOffset) => {
      const targetRow = rowIndex + rowOffset
      const targetCol = colIndex + colOffset

      // 열 추가
      const maxCols = columns.value.length
      if (targetCol >= maxCols) {
        columns.value.push(`Column ${columns.value.length + 1}`)
        rowData.value.forEach((r) => r.push(""))
      }

      // 행 추가
      const maxRows = rowData.value.length
      if (targetRow >= maxRows) {
        rowData.value.push(Array(maxCols).fill(""))
      }

      // 값 설정
      rowData.value[targetRow][targetCol] = value
    })
  })
}

const isSelecting = ref(false)
type Position = { rowIndex: number; colIndex: number }
const startCellPos = ref<Position | null>(null)
const startCellRowIndex = computed(() => startCellPos.value?.rowIndex ?? -1)
const startCellColIndex = computed(() => startCellPos.value?.colIndex ?? -1)
const selectedCells = reactive<Position[]>([])

const startCell = (rowIndex: number, colIndex: number) => {
  isSelecting.value = true
  hasDragged.value = false
  startCellPos.value = { rowIndex, colIndex }
  selectedCells.splice(0) // Clear previous selection
  selectedCells.push({ rowIndex, colIndex })
}

function updateSelection(rowIndex: number, colIndex: number) {
  const startRow = Math.min(startCellRowIndex.value, rowIndex)
  const endRow = Math.max(startCellRowIndex.value, rowIndex)
  const startCol = Math.min(startCellColIndex.value, colIndex)
  const endCol = Math.max(startCellColIndex.value, colIndex)

  selectedCells.splice(0)
  for (let r = startRow; r <= endRow; r++) {
    for (let c = startCol; c <= endCol; c++) {
      selectedCells.push({ rowIndex: r, colIndex: c })
    }
  }
}

const dragCell = (rowIndex: number, colIndex: number) => {
  if (!isSelecting.value) return
  hasDragged.value = true
  updateSelection(rowIndex, colIndex)
}

const endSelection = () => {
  isSelecting.value = false
}

const handleMouseMove = (event: Event) => {
  if (isSelecting.value) {
    event.preventDefault() // Prevent text selection
  }
}

const isSelectedCell = (rowIndex: number, colIndex: number) => {
  return selectedCells.some(
    (cell) => cell.rowIndex === rowIndex && cell.colIndex === colIndex
  )
}
</script>
