pos = {
  // chessboard
  b: [
    ["-", "-", "-", "P", "-", "-", "-", "-"],
    ["-", "p", "-", "-", "-", "P", "-", "-"],
    ["-", "-", "-", "p", "-", "-", "-", "-"],
    ["-", "p", "P", "-", "-", "-", "-", "-"],
    ["-", "-", "p", "P", "-", "-", "-", "-"],
    ["-", "-", "-", "-", "-", "-", "-", "-"],
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
function pawns_mg(pos, square) {
  if (square == null) return sum(pos, pawns_mg);
  var v = 0;
  if (doubled_isolated(pos, square)) {
    v -= 11;
    square.x == 3 && square.y == 2 && console.log("double isolated", v);
  } else if (isolated(pos, square)) {
    v -= 5;
    square.x == 3 && square.y == 2 && console.log("isolated", v);
  } else if (backward(pos, square)) {
    v -= 9;
    square.x == 3 && square.y == 2 && console.log("backward", v);
  }

  v -= doubled(pos, square) * 11;
  square.x == 3 && square.y == 2 && console.log("doubled", v);

  v += connected(pos, square) ? connected_bonus(pos, square) : 0;
  square.x == 3 && square.y == 2 && console.log("connected", v);

  v -= 13 * weak_unopposed_pawn(pos, square);
  square.x == 3 && square.y == 2 && console.log("weak", v);

  v += [0, -11, -3][blocked(pos, square)];
  square.x == 3 && square.y == 2 && console.log("blocked", v);

  console.log(square, v);
  return v;
}
function doubled_isolated(pos, square) {
  if (square == null) return sum(pos, doubled_isolated);
  if (board(pos, square.x, square.y) != "P") return 0;
  if (isolated(pos, square)) {
    var obe = 0,
      eop = 0,
      ene = 0;
    for (var y = 0; y < 8; y++) {
      if (y > square.y && board(pos, square.x, y) == "P") obe++;
      if (y < square.y && board(pos, square.x, y) == "p") eop++;
      if (
        board(pos, square.x - 1, y) == "p" ||
        board(pos, square.x + 1, y) == "p"
      )
        ene++;
    }
    console.log(obe, eop, ene);
    if (obe > 0 && ene == 0 && eop > 0) return 1;
  }
  return 0;
}
function isolated(pos, square) {
  if (square == null) return sum(pos, isolated);
  if (board(pos, square.x, square.y) != "P") return 0;
  for (var y = 0; y < 8; y++) {
    if (board(pos, square.x - 1, y) == "P") return 0;
    if (board(pos, square.x + 1, y) == "P") return 0;
  }
  return 1;
}

function backward(pos, square) {
  if (square == null) return sum(pos, backward);
  if (board(pos, square.x, square.y) != "P") return 0;
  for (var y = square.y; y < 8; y++) {
    if (
      board(pos, square.x - 1, y) == "P" ||
      board(pos, square.x + 1, y) == "P"
    )
      return 0;
  }
  if (
    board(pos, square.x - 1, square.y - 2) == "p" ||
    board(pos, square.x + 1, square.y - 2) == "p" ||
    board(pos, square.x, square.y - 1) == "p"
  )
    return 1;
  return 0;
}

function doubled(pos, square) {
  if (square == null) return sum(pos, doubled);
  if (board(pos, square.x, square.y) != "P") return 0;
  if (board(pos, square.x, square.y + 1) != "P") return 0;
  if (board(pos, square.x - 1, square.y + 1) == "P") return 0;
  if (board(pos, square.x + 1, square.y + 1) == "P") return 0;
  return 1;
}

function connected(pos, square) {
  if (square == null) return sum(pos, connected);
  if (supported(pos, square) || phalanx(pos, square)) return 1;
  return 0;
}
function supported(pos, square) {
  if (square == null) return sum(pos, supported);
  if (board(pos, square.x, square.y) != "P") return 0;
  return (
    (board(pos, square.x - 1, square.y + 1) == "P" ? 1 : 0) +
    (board(pos, square.x + 1, square.y + 1) == "P" ? 1 : 0)
  );
}
function phalanx(pos, square) {
  if (square == null) return sum(pos, phalanx);
  if (board(pos, square.x, square.y) != "P") return 0;
  if (board(pos, square.x - 1, square.y) == "P") return 1;
  if (board(pos, square.x + 1, square.y) == "P") return 1;
  return 0;
}
function connected_bonus(pos, square) {
  if (square == null) return sum(pos, connected_bonus);
  if (!connected(pos, square)) return 0;
  var seed = [0, 7, 8, 12, 29, 48, 86];
  var op = opposed(pos, square);
  var ph = phalanx(pos, square);
  var su = supported(pos, square);
  var r = rank(pos, square);
  if (r < 2 || r > 7) return 0;
  return seed[r - 1] * (2 + ph - op) + 21 * su;
}
function opposed(pos, square) {
  if (square == null) return sum(pos, opposed);
  if (board(pos, square.x, square.y) != "P") return 0;

  for (var y = 0; y < square.y; y++) {
    if (board(pos, square.x, y) == "p") return 1;
  }
  return 0;
}
function rank(pos, square) {
  if (square == null) return sum(pos, rank);
  return 8 - square.y;
}
function weak_unopposed_pawn(pos, square) {
  if (square == null) return sum(pos, weak_unopposed_pawn);
  if (opposed(pos, square)) return 0;
  var v = 0;
  if (isolated(pos, square)) v++;
  else if (backward(pos, square)) v++;
  return v;
}

function blocked(pos, square) {
  if (square == null) return sum(pos, blocked);
  if (board(pos, square.x, square.y) != "P") return 0;
  if (square.y != 2 && square.y != 3) return 0;
  if (board(pos, square.x, square.y - 1) != "p") return 0;
  return 4 - square.y;
}

console.log(pawns_mg(pos));
