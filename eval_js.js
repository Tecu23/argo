pos = {
  // chessboard
  b: [
    ["-", "-", "-", "-", "-", "-", "-", "-"],
    ["-", "-", "-", "-", "-", "-", "-", "-"],
    ["-", "-", "p", "-", "P", "-", "-", "-"],
    ["-", "-", "p", "-", "P", "-", "-", "-"],
    ["-", "-", "p", "-", "P", "B", "-", "-"],
    ["-", "-", "p", "-", "P", "-", "-", "-"],
    ["-", "-", "-", "-", "-", "-", "-", "-"],
    ["-", "-", "-", "-", "-", "-", "-", "-"],
  ],

  // castling rights
  c: [false, false, false, false],

  // enpassant
  e: null,

  // side to move
  w: true,

  // move counts
  m: [0, 1],
};
function board(pos, x, y) {
  if (x >= 0 && x <= 7 && y >= 0 && y <= 7) return pos.b[x][y];
  return "x";
}
function colorflip(pos) {
  var board = new Array(8);
  for (var i = 0; i < 8; i++) board[i] = new Array(8);
  for (x = 0; x < 8; x++)
    for (y = 0; y < 8; y++) {
      board[x][y] = pos.b[x][7 - y];
      var color = board[x][y].toUpperCase() == board[x][y];
      board[x][y] = color
        ? board[x][y].toLowerCase()
        : board[x][y].toUpperCase();
    }
  return {
    b: board,
    c: [pos.c[2], pos.c[3], pos.c[0], pos.c[1]],
    e: pos.e == null ? null : [pos.e[0], 7 - pos.e[1]],
    w: !pos.w,
    m: [pos.m[0], pos.m[1]],
  };
}
function sum(pos, func, param) {
  var sum = 0;
  for (var x = 0; x < 8; x++)
    for (var y = 0; y < 8; y++) sum += func(pos, { x: x, y: y }, param);
  return sum;
}

function bishop_pawns(pos, square) {
  if (square == null) return sum(pos, bishop_pawns);
  if (board(pos, square.x, square.y) != "B") return 0;
  var c = (square.x + square.y) % 2,
    v = 0;
  var blocked = 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y) == "P" && c == (x + y) % 2) {
        v++;
      }
      if (
        board(pos, x, y) == "P" &&
        x > 1 &&
        x < 6 &&
        board(pos, x, y - 1) != "-"
      )
        blocked++;
    }
  }
  console.log(
    v,
    blocked,
    pawn_attack(pos, square) > 0 ? 0 : 1,
    pawn_attack(pos, square),
  );
  return v * (blocked + (pawn_attack(pos, square) > 0 ? 0 : 1));
}

function pawn_attack(pos, square) {
  if (square == null) return sum(pos, pawn_attack);
  console.log(square);
  var v = 0;
  if (board(pos, square.x - 1, square.y + 1) == "P") {
    v++;
  }
  if (board(pos, square.x + 1, square.y + 1) == "P") {
    v++;
  }
  return v;
}

console.log(bishop_pawns(pos));
