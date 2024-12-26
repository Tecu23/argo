pos = {
  // chessboard

  b: [
    ["-", "p", "-", "n", "-", "P", "-", "-"],
    ["-", "-", "p", "-", "-", "-", "P", "B"],
    ["-", "r", "-", "p", "-", "P", "-", "R"],
    ["-", "-", "-", "-", "P", "-", "-", "N"],
    ["-", "-", "p", "-", "-", "P", "-", "-"],
    ["-", "p", "-", "-", "-", "-", "P", "-"],
    ["-", "-", "p", "-", "-", "P", "-", "-"],
    ["-", "k", "-", "-", "-", "-", "K", "-"],
  ],
  // castling rights
  c: [false, false, false, false],

  // enpassant
  e: null,

  // side to move
  w: true,

  // move counts
  m: [0, 17],
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

function king_mg(pos) {
  var v = 0;
  var kd = king_danger(pos);
  v -= shelter_strength(pos);
  v += shelter_storm(pos);
  v += ((kd * kd) / 4096) << 0;
  v += 8 * flank_attack(pos);
  v += 17 * pawnless_flank(pos);
  return v;
}

function king_danger(pos) {
  var count = king_attackers_count(pos);
  var weight = king_attackers_weight(pos);
  var kingAttacks = king_attacks(pos);
  var weak = weak_bonus(pos);
  var unsafeChecks = unsafe_checks(pos);
  var blockersForKing = blockers_for_king(pos);
  var kingFlankAttack = flank_attack(pos);
  var kingFlankDefense = flank_defense(pos);
  var noQueen = queen_count(pos) > 0 ? 0 : 1;
  var v =
    count * weight +
    69 * kingAttacks +
    185 * weak -
    100 * (knight_defender(colorflip(pos)) > 0) +
    148 * unsafeChecks +
    98 * blockersForKing -
    4 * kingFlankDefense +
    (((3 * kingFlankAttack * kingFlankAttack) / 8) << 0) -
    873 * noQueen -
    (((6 * (shelter_strength(pos) - shelter_storm(pos))) / 8) << 0) +
    mobility_mg(pos) -
    mobility_mg(colorflip(pos)) +
    37 +
    ((772 * Math.min(safe_check(pos, null, 3), 1.45)) << 0) +
    ((1084 * Math.min(safe_check(pos, null, 2), 1.75)) << 0) +
    ((645 * Math.min(safe_check(pos, null, 1), 1.5)) << 0) +
    ((792 * Math.min(safe_check(pos, null, 0), 1.62)) << 0);
  if (v > 100) return v;
  return 0;
}

function king_attackers_count(pos, square) {
  if (square == null) return sum(pos, king_attackers_count);
  if ("PNBRQ".indexOf(board(pos, square.x, square.y)) < 0) return 0;
  if (board(pos, square.x, square.y) == "P") {
    var v = 0;
    for (var dir = -1; dir <= 1; dir += 2) {
      var fr = board(pos, square.x + dir * 2, square.y) == "P";
      if (
        square.x + dir >= 0 &&
        square.x + dir <= 7 &&
        king_ring(pos, { x: square.x + dir, y: square.y - 1 }, true)
      )
        v = v + (fr ? 0.5 : 1);
    }
    return v;
  }
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      var s2 = { x: x, y: y };
      if (king_ring(pos, s2)) {
        if (
          knight_attack(pos, s2, square) ||
          bishop_xray_attack(pos, s2, square) ||
          rook_xray_attack(pos, s2, square) ||
          queen_attack(pos, s2, square)
        )
          return 1;
      }
    }
  }
  return 0;
}

function king_ring(pos, square, full) {
  if (square == null) return sum(pos, king_ring);
  if (
    !full &&
    board(pos, square.x + 1, square.y - 1) == "p" &&
    board(pos, square.x - 1, square.y - 1) == "p"
  )
    return 0;
  for (var ix = -2; ix <= 2; ix++) {
    for (var iy = -2; iy <= 2; iy++) {
      if (
        board(pos, square.x + ix, square.y + iy) == "k" &&
        ((ix >= -1 && ix <= 1) || square.x + ix == 0 || square.x + ix == 7) &&
        ((iy >= -1 && iy <= 1) || square.y + iy == 0 || square.y + iy == 7)
      )
        return 1;
    }
  }
  return 0;
}

function knight_attack(pos, square, s2) {
  if (square == null) return sum(pos, knight_attack);
  var v = 0;
  for (var i = 0; i < 8; i++) {
    var ix = ((i > 3) + 1) * ((i % 4 > 1) * 2 - 1);
    var iy = (2 - (i > 3)) * ((i % 2 == 0) * 2 - 1);
    var b = board(pos, square.x + ix, square.y + iy);
    if (
      b == "N" &&
      (s2 == null || (s2.x == square.x + ix && s2.y == square.y + iy)) &&
      !pinned(pos, { x: square.x + ix, y: square.y + iy })
    )
      v++;
  }
  return v;
}

function pinned(pos, square) {
  if (square == null) return sum(pos, pinned);
  if ("PNBRQK".indexOf(board(pos, square.x, square.y)) < 0) return 0;
  return pinned_direction(pos, square) > 0 ? 1 : 0;
}

function pinned_direction(pos, square) {
  if (square == null) return sum(pos, pinned_direction);
  if ("PNBRQK".indexOf(board(pos, square.x, square.y).toUpperCase()) < 0)
    return 0;
  var color = 1;
  if ("PNBRQK".indexOf(board(pos, square.x, square.y)) < 0) color = -1;
  if (square.x == 5 && square.y == 6) {
    console.log(color);
  }
  for (var i = 0; i < 8; i++) {
    var ix = ((i + (i > 3)) % 3) - 1;
    var iy = (((i + (i > 3)) / 3) << 0) - 1;
    var king = false;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (b == "K") king = true;
      if (b != "-") break;
    }
    if (king) {
      for (var d = 1; d < 8; d++) {
        var b = board(pos, square.x - d * ix, square.y - d * iy);
        if (
          b == "q" ||
          (b == "b" && ix * iy != 0) ||
          (b == "r" && ix * iy == 0)
        ) {
          if (square.x == 5 && square.y == 6) {
            console.log(
              b == "q",
              b == "b" && ix * iy != 0,
              b == "r" && ix * iy == 0,
            );
          }
          return Math.abs(ix + iy * 3) * color;
        }
        if (b != "-") break;
      }
    }
  }
  return 0;
}

function bishop_xray_attack(pos, square, s2) {
  if (square == null) return sum(pos, bishop_xray_attack);
  var v = 0;
  for (var i = 0; i < 4; i++) {
    var ix = (i > 1) * 2 - 1;
    var iy = (i % 2 == 0) * 2 - 1;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (
        b == "B" &&
        (s2 == null || (s2.x == square.x + d * ix && s2.y == square.y + d * iy))
      ) {
        var dir = pinned_direction(pos, {
          x: square.x + d * ix,
          y: square.y + d * iy,
        });
        if (dir == 0 || Math.abs(ix + iy * 3) == dir) v++;
      }
      if (b != "-" && b != "Q" && b != "q") break;
    }
  }
  return v;
}

function rook_xray_attack(pos, square, s2) {
  if (square == null) return sum(pos, rook_xray_attack);
  var v = 0;
  for (var i = 0; i < 4; i++) {
    var ix = i == 0 ? -1 : i == 1 ? 1 : 0;
    var iy = i == 2 ? -1 : i == 3 ? 1 : 0;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (
        b == "R" &&
        (s2 == null || (s2.x == square.x + d * ix && s2.y == square.y + d * iy))
      ) {
        var dir = pinned_direction(pos, {
          x: square.x + d * ix,
          y: square.y + d * iy,
        });
        if (dir == 0 || Math.abs(ix + iy * 3) == dir) v++;
      }
      if (b != "-" && b != "R" && b != "Q" && b != "q") break;
    }
  }

  return v;
}

function queen_attack(pos, square, s2) {
  if (square == null) return sum(pos, queen_attack);
  var v = 0;
  for (var i = 0; i < 8; i++) {
    var ix = ((i + (i > 3)) % 3) - 1;
    var iy = (((i + (i > 3)) / 3) << 0) - 1;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (
        b == "Q" &&
        (s2 == null || (s2.x == square.x + d * ix && s2.y == square.y + d * iy))
      ) {
        var dir = pinned_direction(pos, {
          x: square.x + d * ix,
          y: square.y + d * iy,
        });
        if (dir == 0 || Math.abs(ix + iy * 3) == dir) v++;
      }
      if (b != "-") break;
    }
  }
  return v;
}

function king_attackers_weight(pos, square) {
  if (square == null) return sum(pos, king_attackers_weight);
  if (king_attackers_count(pos, square)) {
    return [0, 81, 52, 44, 10]["PNBRQ".indexOf(board(pos, square.x, square.y))];
  }
  return 0;
}

function king_attacks(pos, square) {
  if (square == null) return sum(pos, king_attacks);
  if ("NBRQ".indexOf(board(pos, square.x, square.y)) < 0) return 0;
  if (king_attackers_count(pos, square) == 0) return 0;
  var kx = 0,
    ky = 0,
    v = 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y) == "k") {
        kx = x;
        ky = y;
      }
    }
  }
  for (var x = kx - 1; x <= kx + 1; x++) {
    for (var y = ky - 1; y <= ky + 1; y++) {
      var s2 = { x: x, y: y };
      if (x >= 0 && y >= 0 && x <= 7 && y <= 7 && (x != kx || y != ky)) {
        v += knight_attack(pos, s2, square);
        v += bishop_xray_attack(pos, s2, square);
        v += rook_xray_attack(pos, s2, square);
        v += queen_attack(pos, s2, square);
      }
    }
  }
  return v;
}

function weak_bonus(pos, square) {
  if (square == null) return sum(pos, weak_bonus);
  if (!weak_squares(pos, square)) return 0;
  if (!king_ring(pos, square)) return 0;
  return 1;
}

function weak_squares(pos, square) {
  if (square == null) return sum(pos, weak_squares);
  if (attack(pos, square)) {
    var pos2 = colorflip(pos);
    var attack = attack(pos2, { x: square.x, y: 7 - square.y });
    if (attack >= 2) return 0;
    if (attack == 0) return 1;
    if (
      king_attack(pos2, { x: square.x, y: 7 - square.y }) ||
      queen_attack(pos2, { x: square.x, y: 7 - square.y })
    )
      return 1;
  }
  return 0;
}

function attack(pos, square) {
  if (square == null) return sum(pos, attack);
  var v = 0;
  v += pawn_attack(pos, square);
  if (square.x == 1 && square.y == 1) {
    console.log(v);
  }
  v += king_attack(pos, square);
  if (square.x == 1 && square.y == 1) {
    console.log(v);
  }
  v += knight_attack(pos, square);
  if (square.x == 1 && square.y == 1) {
    console.log(v);
  }
  v += bishop_xray_attack(pos, square);
  if (square.x == 1 && square.y == 1) {
    console.log(v);
  }
  v += rook_xray_attack(pos, square);
  if (square.x == 1 && square.y == 1) {
    console.log(v);
  }
  v += queen_attack(pos, square);
  if (square.x == 1 && square.y == 1) {
    console.log(v);
  }
  return v;
}

function pawn_attack(pos, square) {
  if (square == null) return sum(pos, pawn_attack);
  var v = 0;
  if (board(pos, square.x - 1, square.y + 1) == "P") v++;
  if (board(pos, square.x + 1, square.y + 1) == "P") v++;
  return v;
}

function king_attack(pos, square) {
  if (square == null) return sum(pos, king_attack);
  for (var i = 0; i < 8; i++) {
    var ix = ((i + (i > 3)) % 3) - 1;
    var iy = (((i + (i > 3)) / 3) << 0) - 1;
    if (board(pos, square.x + ix, square.y + iy) == "K") return 1;
  }

  return 0;
}

function unsafe_checks(pos, square) {
  if (square == null) return sum(pos, unsafe_checks);
  if (check(pos, square, 0) && safe_check(pos, null, 0) == 0) return 1;
  if (check(pos, square, 1) && safe_check(pos, null, 1) == 0) return 1;
  if (check(pos, square, 2) && safe_check(pos, null, 2) == 0) return 1;
  return 0;
}

function check(pos, square, type) {
  if (square == null) return sum(pos, check);
  if (
    (rook_xray_attack(pos, square) &&
      (type == null || type == 2 || type == 4)) ||
    (queen_attack(pos, square) && (type == null || type == 3))
  ) {
    for (var i = 0; i < 4; i++) {
      var ix = i == 0 ? -1 : i == 1 ? 1 : 0;
      var iy = i == 2 ? -1 : i == 3 ? 1 : 0;
      for (var d = 1; d < 8; d++) {
        var b = board(pos, square.x + d * ix, square.y + d * iy);
        if (b == "k") return 1;
        if (b != "-" && b != "q") break;
      }
    }
  }
  if (
    (bishop_xray_attack(pos, square) &&
      (type == null || type == 1 || type == 4)) ||
    (queen_attack(pos, square) && (type == null || type == 3))
  ) {
    for (var i = 0; i < 4; i++) {
      var ix = (i > 1) * 2 - 1;
      var iy = (i % 2 == 0) * 2 - 1;
      for (var d = 1; d < 8; d++) {
        var b = board(pos, square.x + d * ix, square.y + d * iy);
        if (b == "k") return 1;
        if (b != "-" && b != "q") break;
      }
    }
  }
  if (knight_attack(pos, square) && (type == null || type == 0 || type == 4)) {
    if (
      board(pos, square.x + 2, square.y + 1) == "k" ||
      board(pos, square.x + 2, square.y - 1) == "k" ||
      board(pos, square.x + 1, square.y + 2) == "k" ||
      board(pos, square.x + 1, square.y - 2) == "k" ||
      board(pos, square.x - 2, square.y + 1) == "k" ||
      board(pos, square.x - 2, square.y - 1) == "k" ||
      board(pos, square.x - 1, square.y + 2) == "k" ||
      board(pos, square.x - 1, square.y - 2) == "k"
    )
      return 1;
  }
  return 0;
}

function safe_check(pos, square, type) {
  if (square == null) return sum(pos, safe_check, type);
  if ("PNBRQK".indexOf(board(pos, square.x, square.y)) >= 0) return 0;
  if (!check(pos, square, type)) return 0;
  var pos2 = colorflip(pos);
  if (type == 3 && safe_check(pos, square, 2)) return 0;
  if (type == 1 && safe_check(pos, square, 3)) return 0;
  if (
    (!attack(pos2, { x: square.x, y: 7 - square.y }) ||
      (weak_squares(pos, square) && attack(pos, square) > 1)) &&
    (type != 3 || !queen_attack(pos2, { x: square.x, y: 7 - square.y }))
  )
    return 1;
  return 0;
}
function blockers_for_king(pos, square) {
  if (square == null) return sum(pos, blockers_for_king);
  if (pinned_direction(colorflip(pos), { x: square.x, y: 7 - square.y })) {
    console.log(square);
    return 1;
  }
  return 0;
}
function flank_attack(pos, square) {
  if (square == null) return sum(pos, flank_attack);
  if (square.y > 4) return 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y) == "k") {
        if (x == 0 && square.x > 2) return 0;
        if (x < 3 && square.x > 3) return 0;
        if (x >= 3 && x < 5 && (square.x < 2 || square.x > 5)) return 0;
        if (x >= 5 && square.x < 4) return 0;
        if (x == 7 && square.x < 5) return 0;
      }
    }
  }
  var a = attack(pos, square);
  if (!a) return 0;
  return a > 1 ? 2 : 1;
}

function flank_defense(pos, square) {
  if (square == null) return sum(pos, flank_defense);
  if (square.y > 4) return 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y) == "k") {
        if (x == 0 && square.x > 2) return 0;
        if (x < 3 && square.x > 3) return 0;
        if (x >= 3 && x < 5 && (square.x < 2 || square.x > 5)) return 0;
        if (x >= 5 && square.x < 4) return 0;
        if (x == 7 && square.x < 5) return 0;
      }
    }
  }
  return attack(colorflip(pos), { x: square.x, y: 7 - square.y }) > 0 ? 1 : 0;
}

function queen_count(pos, square) {
  if (square == null) return sum(pos, queen_count);
  if (board(pos, square.x, square.y) == "Q") return 1;
  return 0;
}

function knight_defender(pos, square) {
  if (square == null) return sum(pos, knight_defender);
  if (knight_attack(pos, square) && king_attack(pos, square)) return 1;
  return 0;
}

function shelter_strength(pos, square) {
  var w = 0,
    s = 1024,
    tx = null;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (
        board(pos, x, y) == "k" ||
        (pos.c[2] && x == 6 && y == 0) ||
        (pos.c[3] && x == 2 && y == 0)
      ) {
        var w1 = strength_square(pos, { x: x, y: y });
        console.log(x, y, w1);
        var s1 = storm_square(pos, { x: x, y: y });
        if (s1 - w1 < s - w) {
          w = w1;
          s = s1;
          tx = Math.max(1, Math.min(6, x));
        }
      }
    }
  }
  if (square == null) return w;
  if (
    tx != null &&
    board(pos, square.x, square.y) == "p" &&
    square.x >= tx - 1 &&
    square.x <= tx + 1
  ) {
    for (var y = square.y - 1; y >= 0; y--)
      if (board(pos, square.x, y) == "p") return 0;
    return 1;
  }
  return 0;
}
function strength_square(pos, square) {
  if (square == null) return sum(pos, strength_square);
  var v = 5;
  var kx = Math.min(6, Math.max(1, square.x));
  var weakness = [
    [-6, 81, 93, 58, 39, 18, 25],
    [-43, 61, 35, -49, -29, -11, -63],
    [-10, 75, 23, -2, 32, 3, -45],
    [-39, -13, -29, -52, -48, -67, -166],
  ];
  for (var x = kx - 1; x <= kx + 1; x++) {
    var us = 0;
    for (var y = 7; y >= square.y; y--) {
      if (
        board(pos, x, y) == "p" &&
        board(pos, x - 1, y + 1) != "P" &&
        board(pos, x + 1, y + 1) != "P"
      )
        us = y;
    }
    var f = Math.min(x, 7 - x);
    console.log(weakness[f][us]);
    v += weakness[f][us] || 0;
  }
  return v;
}
function shelter_storm(pos, square) {
  var w = 0,
    s = 1024,
    tx = null;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (
        board(pos, x, y) == "k" ||
        (pos.c[2] && x == 6 && y == 0) ||
        (pos.c[3] && x == 2 && y == 0)
      ) {
        var w1 = strength_square(pos, { x: x, y: y });
        var s1 = storm_square(pos, { x: x, y: y });
        if (s1 - w1 < s - w) {
          w = w1;
          s = s1;
          tx = Math.max(1, Math.min(6, x));
        }
      }
    }
  }
  if (square == null) return s;
  if (
    tx != null &&
    board(pos, square.x, square.y).toUpperCase() == "P" &&
    square.x >= tx - 1 &&
    square.x <= tx + 1
  ) {
    for (var y = square.y - 1; y >= 0; y--)
      if (board(pos, square.x, y) == board(pos, square.x, square.y)) return 0;
    return 1;
  }
  return 0;
}

function storm_square(pos, square, eg) {
  if (square == null) return sum(pos, storm_square);
  var v = 0,
    ev = 5;
  var kx = Math.min(6, Math.max(1, square.x));
  var unblockedstorm = [
    [85, -289, -166, 97, 50, 45, 50],
    [46, -25, 122, 45, 37, -10, 20],
    [-6, 51, 168, 34, -2, -22, -14],
    [-15, -11, 101, 4, 11, -15, -29],
  ];
  var blockedstorm = [
    [0, 0, 76, -10, -7, -4, -1],
    [0, 0, 78, 15, 10, 6, 2],
  ];
  for (var x = kx - 1; x <= kx + 1; x++) {
    var us = 0,
      them = 0;
    for (var y = 7; y >= square.y; y--) {
      if (
        board(pos, x, y) == "p" &&
        board(pos, x - 1, y + 1) != "P" &&
        board(pos, x + 1, y + 1) != "P"
      )
        us = y;
      if (board(pos, x, y) == "P") them = y;
    }
    var f = Math.min(x, 7 - x);
    if (us > 0 && them == us + 1) {
      v += blockedstorm[0][them];
      ev += blockedstorm[1][them];
    } else v += unblockedstorm[f][them];
  }
  return eg ? ev : v;
}

function mobility_mg(pos, square) {
  if (square == null) return sum(pos, mobility_mg);
  return mobility_bonus(pos, square, true);
}

function mobility_bonus(pos, square, mg) {
  if (square == null) return sum(pos, mobility_bonus, mg);
  var bonus = mg
    ? [
        [-62, -53, -12, -4, 3, 13, 22, 28, 33],
        [-48, -20, 16, 26, 38, 51, 55, 63, 63, 68, 81, 81, 91, 98],
        [-60, -20, 2, 3, 3, 11, 22, 31, 40, 40, 41, 48, 57, 57, 62],
        [
          -30, -12, -8, -9, 20, 23, 23, 35, 38, 53, 64, 65, 65, 66, 67, 67, 72,
          72, 77, 79, 93, 108, 108, 108, 110, 114, 114, 116,
        ],
      ]
    : [
        [-81, -56, -31, -16, 5, 11, 17, 20, 25],
        [-59, -23, -3, 13, 24, 42, 54, 57, 65, 73, 78, 86, 88, 97],
        [-78, -17, 23, 39, 70, 99, 103, 121, 134, 139, 158, 164, 168, 169, 172],
        [
          -48, -30, -7, 19, 40, 55, 59, 75, 78, 96, 96, 100, 121, 127, 131, 133,
          136, 141, 147, 150, 151, 168, 168, 171, 182, 182, 192, 219,
        ],
      ];
  var i = "NBRQ".indexOf(board(pos, square.x, square.y));
  if (i < 0) return 0;
  return bonus[i][mobility(pos, square)];
}

function mobility(pos, square) {
  if (square == null) return sum(pos, mobility);
  var v = 0;
  var b = board(pos, square.x, square.y);
  if ("NBRQ".indexOf(b) < 0) return 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      var s2 = { x: x, y: y };
      if (!mobility_area(pos, s2)) continue;
      if (b == "N" && knight_attack(pos, s2, square) && board(pos, x, y) != "Q")
        v++;
      if (
        b == "B" &&
        bishop_xray_attack(pos, s2, square) &&
        board(pos, x, y) != "Q"
      )
        v++;
      if (b == "R" && rook_xray_attack(pos, s2, square)) v++;
      if (b == "Q" && queen_attack(pos, s2, square)) v++;
    }
  }
  return v;
}

function mobility_area(pos, square) {
  if (square == null) return sum(pos, mobility_area);
  if (board(pos, square.x, square.y) == "K") return 0;
  if (board(pos, square.x, square.y) == "Q") return 0;
  if (board(pos, square.x - 1, square.y - 1) == "p") return 0;
  if (board(pos, square.x + 1, square.y - 1) == "p") return 0;
  if (
    board(pos, square.x, square.y) == "P" &&
    (rank(pos, square) < 4 || board(pos, square.x, square.y - 1) != "-")
  )
    return 0;
  if (blockers_for_king(colorflip(pos), { x: square.x, y: 7 - square.y }))
    return 0;
  return 1;
}

function rank(pos, square) {
  if (square == null) return sum(pos, rank);
  return 8 - square.y;
}

function pawnless_flank(pos) {
  var pawns = [0, 0, 0, 0, 0, 0, 0, 0],
    kx = 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y).toUpperCase() == "P") pawns[x]++;
      if (board(pos, x, y) == "k") kx = x;
    }
  }
  var sum;
  if (kx == 0) sum = pawns[0] + pawns[1] + pawns[2];
  else if (kx < 3) sum = pawns[0] + pawns[1] + pawns[2] + pawns[3];
  else if (kx < 5) sum = pawns[2] + pawns[3] + pawns[4] + pawns[5];
  else if (kx < 7) sum = pawns[4] + pawns[5] + pawns[6] + pawns[7];
  else sum = pawns[5] + pawns[6] + pawns[7];
  return sum == 0 ? 1 : 0;
}

function restricted(pos, square) {
  if (square == null) return sum(pos, restricted);

  if (attack(pos, square) == 0) {
    console.log("1");
    return 0;
  }

  var pos2 = colorflip(pos);

  if (!attack(pos2, { x: square.x, y: 7 - square.y })) {
    console.log("2");
    return 0;
  }

  if (pawn_attack(pos2, { x: square.x, y: 7 - square.y }) > 0) {
    console.log("3");
    return 0;
  }
  if (
    attack(pos2, { x: square.x, y: 7 - square.y }) > 1 &&
    attack(pos, square) == 1
  ) {
    console.log("4");
    return 0;
  }
  return 1;
}

console.log(restricted(pos, { x: 1, y: 6 }));
