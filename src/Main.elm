module Main exposing (..)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Css exposing (..)
import Browser

--- Main

main = Browser.sandbox {
    init = init,
    update = update,
    view = view
    }

-- Model
type Msg index value = 
    Set index value

type alias Model = {
        cells: List (Maybe String),
        actualWords: List String
    }

init : Model
init = {
        cells = List.map (\_ -> Nothing) (List.range 1 10),
        actualWords = ["cat", "cup", "ted", "ped", "top"]
    }

-- Update

updateSingleElementInList : List a -> Int -> a -> List a
updateSingleElementInList list pos newValue =
    (List.take pos list) ++ [newValue] ++ (List.drop (pos + 1) list)

update : Msg Int String -> Model -> Model
update msg model = 
    case msg of
        Set index value -> 
            let
                _ = Debug.log "hello world" (index)
            in
                { model | cells = (updateSingleElementInList model.cells index (Just value)) }
            

--- View
applicationContainerStyle : List (Attribute msg)
applicationContainerStyle =
    [
        style "background-color" "#033860",
        style "display" "grid",
        style "grid-template-rows" "100%",
        style "grid-template-columns" "15% 70% 15%",
        style "min-height" "100vh"
    ]

playerColumn : String -> List (Attribute msg)
playerColumn color =
    [
        style "background-color" color,
        style "min-height" "100vh",
        style "text-align" "center"
    ]

boardRow : List (Attribute msg)
boardRow = 
    [
        -- style "background-color" "white",
        style "padding" "2%",
        style "display" "flex",
        style "justify-content" "center"
    ]

boardCell : List (Attribute msg)
boardCell = 
    [
    style "width" "100px",
    style "height" "100px",
    style "outline" "none",
    style "border" "none",
    style "background-color" "white",
    style "margin" "2%",
    style "font-size" "2.5em",
    style "display" "flex",
    style "justify-content" "center",
    style "align-items" "center",
    style "cursor" "pointer"
    ]

boardWrapper : List (Attribute msg)
boardWrapper =
    [
        style "background-color" "#033860",
        style "display" "flex",
        style "flex-direction" "column",
        style "justify-content" "center",
        style "align-items" "center"
    ]

board : List (Attribute msg)
board = 
    [
        style "background-color" "#033860",
        style "display" "grid",
        style "grid-template-rows" "33% 33% 33%",
        style "width" "420px"
    ]

getValueAt : Int -> List a -> Maybe a
getValueAt index list =
    List.head (List.drop index list)

getCellValue : Int ->Model -> String
getCellValue index model = 
    case getValueAt index model.cells of
        Just (Just value) -> value
        _ -> ""

generateSingleCell : Int -> Model -> Html (Msg Int String)
generateSingleCell index model =
    div 
        (boardCell ++ [ onClick (Set index (String.fromInt index))]) 
        [text (getCellValue index model)]

generateSingleRow : Int -> Int -> Model -> Html (Msg Int String)
generateSingleRow rowIndex size model =
    div boardRow 
        (List.map 
            (\count -> (generateSingleCell ((rowIndex - 1) * size + count - 1) model)) 
            (List.range 1 size))

generate : Int -> Model -> List (Html (Msg Int String))
generate size model =
    List.map (\index -> generateSingleRow index size model) (List.range 1 size)

middleColumn : List (Attribute msg)
middleColumn = 
    [
        style "display" "flex",
        style "flex-direction" "column",
        style "justify-content" "center"
    ]

view : Model -> Html (Msg Int String)
view model = 
    div applicationContainerStyle 
        [
        div (playerColumn "red") [ text "Player 1" ],
        div middleColumn [
            -- board wrapper
            div boardWrapper [
                div [] (generate 3 model)
            ]
        ],
        div (playerColumn "green") [ text "Player 2" ]
        
        ]