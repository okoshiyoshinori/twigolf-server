package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/okoshiyoshinori/twigolf-server/model"
	"github.com/okoshiyoshinori/twigolf-server/util"
)


var style = `{"border":[
        {"type":"bottom","color":"000000","style":1},
        {"type":"top","color":"000000","style":1},
        {"type":"left","color":"000000","style":1},
        {"type":"right","color":"000000","style":1}
        ],
        "alignment": {
          "horizontal": "center",
          "vertical":"center",
          "shrink_to_fit": true
        }}
        `

var styleHeader = `{"border":[
        {"type":"bottom","color":"000000","style":1},
        {"type":"top","color":"000000","style":1},
        {"type":"left","color":"000000","style":1},
        {"type":"right","color":"000000","style":1}
        ],
        "fill": {
          "type":"pattern","color":["#dddddd"],"pattern":1
        },
        "alignment": {
          "horizontal": "center",
          "vertical":"center",
          "shrink_to_fit": true
        }
    }`

type ColumnPositions struct {
  title string 
  date string
  leader string
  place string
}

type ColumnSize struct {
  width int
  height int
}

type ParticipantsRowInfo struct {
  snsid_start_col string
  snsid_size ColumnSize
  name_start_col string
  name_size ColumnSize
  kana_start_col string
  kana_size ColumnSize
  sex_start_col string
  sex_size ColumnSize
  old_start_col string
  old_size ColumnSize
}

type CombinationRowInfo struct {
  stime_start_col string
  stime_size ColumnSize
  in_out_start_col string
  in_out_size ColumnSize
  member1_start_col string
  member1_size ColumnSize
  member2_start_col string
  member2_size ColumnSize
  member3_start_col string
  member3_size ColumnSize
  member4_start_col string
  member4_size ColumnSize
}

var combinationRow = CombinationRowInfo{
  stime_start_col:"C",
  stime_size: ColumnSize{width:2,height:2},
  in_out_start_col:"E",
  in_out_size: ColumnSize{width:2,height:2},
  member1_start_col:"G",
  member1_size: ColumnSize{width:5,height:2},
  member2_start_col:"L",
  member2_size: ColumnSize{width:5,height:2},
  member3_start_col:"Q",
  member3_size: ColumnSize{width:5,height:2},
  member4_start_col:"V",
  member4_size: ColumnSize{width:5,height:2},
}

var position = settingPosition()

var participantsRow = ParticipantsRowInfo{
  snsid_start_col: "C",
  snsid_size: ColumnSize{width:4,height:2},
  name_start_col: "G",
  name_size: ColumnSize{width:8,height:2},
  kana_start_col: "O",
  kana_size: ColumnSize{width:8,height:2},
  sex_start_col: "W",
  sex_size: ColumnSize{width:2,height:2},
  old_start_col: "Y",
  old_size: ColumnSize{width:2,height:2},
}

type Excel struct {
  Fd *excelize.File
  Line int
  participants []*model.Participant
  competition *model.Competition 
  combination []*model.Combination
}

func NewExcel(f *excelize.File,p []*model.Participant,comp *model.Competition,comb []*model.Combination) *Excel {
  return &Excel{
    Fd:f,
    Line: 10,
    participants:p,
    competition:comp,
    combination:comb,
  }
}

func settingPosition() ColumnPositions {
  return ColumnPositions {
    title: "C2",
    date: "G5",
    leader: "S5",
    place:"G7",
  } 
}

func GetCellIntName(col int,row int) string {
  col_name,_ := excelize.ColumnNumberToName(col)
  cell_name,_ := excelize.JoinCellName(col_name,row)
  return cell_name
}

func GetCellStrName(col string,row int) string {
  cell_name,_ := excelize.JoinCellName(col,row)
  return cell_name
}

func (e *Excel) MergeCell(activeSheet string,start_col string,size ColumnSize,header bool) {
  styleID,_ := e.Fd.NewStyle(style)
  styleH,_ := e.Fd.NewStyle(styleHeader)
  num,_ := excelize.ColumnNameToNumber(start_col)
  s_col_num := num
  s_row_num := e.Line 
  end_col_num := s_col_num + size.width - 1
  end_row_num := e.Line + size.height - 1
  s_cell := GetCellIntName(s_col_num,s_row_num)
  e_cell := GetCellIntName(end_col_num,end_row_num)
  e.Fd.MergeCell(activeSheet,s_cell,e_cell)
  e.Fd.SetCellStyle(activeSheet,s_cell,e_cell,styleID)
  if header {
    e.Fd.SetCellStyle(activeSheet,s_cell,e_cell,styleH)
  } else {
    e.Fd.SetCellStyle(activeSheet,s_cell,e_cell,styleID)
  }
}

func (e *Excel) makePatiRow(activeSheet string) {
  e.MergeCell(activeSheet,participantsRow.snsid_start_col,participantsRow.snsid_size,false)
  e.MergeCell(activeSheet,participantsRow.name_start_col,participantsRow.name_size,false)
  e.MergeCell(activeSheet,participantsRow.kana_start_col,participantsRow.kana_size,false)
  e.MergeCell(activeSheet,participantsRow.sex_start_col,participantsRow.sex_size,false)
  e.MergeCell(activeSheet,participantsRow.old_start_col,participantsRow.old_size,false)
}

func (e *Excel) makePartiHeader(activeSheet string) {
  e.MergeCell(activeSheet,participantsRow.snsid_start_col,participantsRow.snsid_size,true)
  e.MergeCell(activeSheet,participantsRow.name_start_col,participantsRow.name_size,true)
  e.MergeCell(activeSheet,participantsRow.kana_start_col,participantsRow.kana_size,true)
  e.MergeCell(activeSheet,participantsRow.sex_start_col,participantsRow.sex_size,true)
  e.MergeCell(activeSheet,participantsRow.old_start_col,participantsRow.old_size,true)
  
  e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.snsid_start_col,e.Line),"SNSID")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.name_start_col,e.Line),"氏名")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.kana_start_col,e.Line),"ふりがな")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.sex_start_col,e.Line),"性別")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.old_start_col,e.Line),"年齢")
}

func (e *Excel) makeCombiRow(activeSheet string) {
  e.MergeCell(activeSheet,combinationRow.stime_start_col,combinationRow.stime_size,false)
  e.MergeCell(activeSheet,combinationRow.in_out_start_col,combinationRow.in_out_size,false)
  e.MergeCell(activeSheet,combinationRow.member1_start_col,combinationRow.member1_size,false)
  e.MergeCell(activeSheet,combinationRow.member2_start_col,combinationRow.member2_size,false)
  e.MergeCell(activeSheet,combinationRow.member3_start_col,combinationRow.member3_size,false)
  e.MergeCell(activeSheet,combinationRow.member4_start_col,combinationRow.member4_size,false)
}

func (e *Excel) makeCombiHeader(activeSheet string) {
  e.MergeCell(activeSheet,combinationRow.stime_start_col,combinationRow.stime_size,true)
  e.MergeCell(activeSheet,combinationRow.in_out_start_col,combinationRow.in_out_size,true)
  e.MergeCell(activeSheet,combinationRow.member1_start_col,combinationRow.member1_size,true)
  e.MergeCell(activeSheet,combinationRow.member2_start_col,combinationRow.member2_size,true)
  e.MergeCell(activeSheet,combinationRow.member3_start_col,combinationRow.member3_size,true)
  e.MergeCell(activeSheet,combinationRow.member4_start_col,combinationRow.member4_size,true)
  
  e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.stime_start_col,e.Line),"時間")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.in_out_start_col,e.Line),"IN/OUT")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member1_start_col,e.Line),"PLAYER")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member2_start_col,e.Line),"PLAYER")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member3_start_col,e.Line),"PLAYER")
  e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member4_start_col,e.Line),"PLAYER")
}

func (e *Excel) Make(activeSheet string) error {
  e.Fd.SetCellStr(activeSheet,position.title,e.competition.Title)
  if e.competition.EventDay != nil {
    e.Fd.SetCellStr(activeSheet,position.date,util.DateToString(e.competition.EventDay))
  }
  if e.competition.User.RealName != nil {
    e.Fd.SetCellStr(activeSheet,position.leader,*e.competition.User.RealName)
  } else {
    e.Fd.SetCellStr(activeSheet,position.leader,e.competition.User.ScreenName)
  }
  if e.competition.PlaceText != nil {
    e.Fd.SetCellStr(activeSheet,position.place,*(e.competition.PlaceText))
  }
  //参加者ヘッダー
  e.makePartiHeader(activeSheet)
  e.Line += participantsRow.snsid_size.height
  //参加者データ
  for _,v := range e.participants {
    e.makePatiRow(activeSheet)
    e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.snsid_start_col,e.Line),v.User.ScreenName)
    if v.User.RealName != nil {
      e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.name_start_col,e.Line),*v.User.RealName)
    }
    if v.User.RealNameKana != nil {
      e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.kana_start_col,e.Line),*v.User.RealNameKana)
    }
    if v.User.Sex != nil {
      e.Fd.SetCellStr(activeSheet,GetCellStrName(participantsRow.sex_start_col,e.Line),util.GetSexToString(*v.User.Sex))
    }
    if v.User.Birthday != nil {
      e.Fd.SetCellInt(activeSheet,GetCellStrName(participantsRow.old_start_col,e.Line),util.CalcAge(*v.User.Birthday))
    }
    e.Line += participantsRow.snsid_size.height
  }

  e.Line += 1 

  //組み合わせヘッダー
  e.makeCombiHeader(activeSheet)
  e.Line += combinationRow.stime_size.height
  for _,v := range e.combination {
    e.makeCombiRow(activeSheet)
    e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.stime_start_col,e.Line),util.TimeToString(v.StartTime))
    e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.in_out_start_col,e.Line),util.GetInOutToString(v.StartInOut))
    if v.Member1 != nil {
      e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member1_start_col,e.Line),util.GetUserRealName(e.participants,*v.Member1))
    }
    if v.Member2 != nil {
      e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member2_start_col,e.Line),util.GetUserRealName(e.participants,*v.Member2))
    }
    if v.Member3 != nil {
      e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member3_start_col,e.Line),util.GetUserRealName(e.participants,*v.Member3))
    }
    if v.Member4 != nil {
      e.Fd.SetCellStr(activeSheet,GetCellStrName(combinationRow.member4_start_col,e.Line),util.GetUserRealName(e.participants,*v.Member4))
    }
    e.Line += combinationRow.stime_size.height
  }
   return nil
}
  

